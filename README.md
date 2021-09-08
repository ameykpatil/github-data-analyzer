# Github Data Analyzer
The application analyzes given github events 

* [Summary](#summary)
* [Setup & Run](#setup--run)
  * [Using application binary](#using-the-application-binary-built-with-makefile)
  * [Using docker container](#using-a-docker-container-image-built-using-the-makefile)
  * [Running Tests](#run-the-tests)
  * [Running Lint](#run-the-linter)
* [Run With Custom Input](#run-with-custom-input)
* [Application Design](#application-design)
* [Possible Enhancements](#possible-enhancements)   

## Summary
Given the directory path to csv files containing github event data, this CLI application can output -

- Top 10 active users sorted by amount of PRs created and commits pushed
- Top 10 repositories sorted by amount of commits pushed
- Top 10 repositories sorted by amount of watch events

This application is written in Golang. Any type of database or data processing engines (such as Apache Spark) have no been used.

The code is written keeping in mind following qualities: 
- Readable 
- Testable 
- Maintainable 
- Extensible 
- Modifiable 
- Performant

Apart from this -
- Basic unit tests are written to demonstrate testing  
- Structured & meaningful commits have been made to be useful while reviewing & also for future references
- Setup section in Readme is provided for all the setup & run instructions

## Setup & Run
Clone the repository
git clone https://github.com/ameykpatil/github-data-analyzer

The application can be run in 2 ways:
1. Using the application binary built with `Makefile`.
2. Using a Docker container image built using the `Makefile`.

#### Using the application binary built with `Makefile`

1. Ensure you are in the root directory of the repository.
2. Run `make build`, it will create the binary `github-data-analyzer` in the root directory of the repository.
3. Run the following command by providing path of the directory where CSV files reside:
```bash
 ./github-data-analyzer all -p=./data/given-data 
```

#### Using a Docker container image built using the `Makefile`.
1. Ensure you are in the root directory of the repository.
2. Run `make docker-build`, this will create the container image `github-data-analyzer:latest`
3. Attach a Docker volume for the directory so that it can be supplies as `-p` (path) value.
4. Run the following by mounting a directory & specifying it in path:
```bash
docker run -v $PWD/data/given-data:/data github-data-analyzer all -p=/data
```
_Note : For other commands & custom inputs see [Run with custom input](#run-with-custom-input) section_

#### Run the tests
Unit tests can be run using `make`
```bash
make test
```

#### Run the linter
Linter can be run using `make` 
But `golangci-lint` will be needed to run it
```bash
make lint
```

## Run With Custom Input

The application is written to be highly Extensible & Flexible. Hence, the custom input is supported in various ways.

- `all` command  
This command serves the purpose of providing the output as specified in the requirements in the summary.  
But it is still possible to provide `limit` (`-l`) flag & change the result from `top 10` to may be `top 15`.  
Following are some examples
```bash
docker run -v $PWD/data/given-data:/data github-data-analyzer all -p=/data -l=15
docker run -v $PWD/data/given-data:/data github-data-analyzer all -p=/data -l=20
```

- `repo` command  
This command serves the purpose of providing the output for top repos.  
It is possible to provide `limit` (`-l`) flag & change the result from `top 10` to may be `top 15`.  
It is also possible to provide `sort` (`-s`) flag to give the specific sorting field based on which top users should be found out.
The valid values for the sort fields are `Commits` or any type of `Event` e.g. `PullRequestEvent`, `WatchEvent` etc.  
Following are some examples
```bash
docker run -v $PWD/data/given-data:/data github-data-analyzer repos -p=/data -l=10 -s=WatchEvent
docker run -v $PWD/data/given-data:/data github-data-analyzer repos -p=/data -l=20 -s=Commits
```

- `users` command  
This command serves the purpose of providing the output for top users.  
It is possible to provide `limit` (`-l`) flag & change the result from `top 10` to may be `top 15`.  
It is also possible to provide `sort` (`-s`) flag to give the list of fields based on which top users should be found out.
The application honors the order of the sort fields provided & consider the sorting in that specific order.
The valid values for the sort fields are `Commits` or any type of `Event` e.g. `PullRequestEvent`, `WatchEvent` etc.  
Following are some examples
```bash
docker run -v $PWD/data/given-data:/data github-data-analyzer users -p=/data -l=15 -s=Commits,PullRequestEvent
docker run -v $PWD/data/given-data:/data github-data-analyzer users -p=/data -l=20 -s=WatchEvent
docker run -v $PWD/data/given-data:/data github-data-analyzer users -p=/data -l=20 -s=ForkEvent,PushEvent
docker run -v $PWD/data/given-data:/data github-data-analyzer users -p=/data -l=20 -s=PullRequestEvent,Commits,PushEvent
```   

## Application Design

- Application has been designed & structured in a layered format. Following diagram should help to visualise the four main layers.  
<img width="400" alt="application layers" src="https://user-images.githubusercontent.com/3050421/132500668-a509ab00-846f-4f89-b459-940122918a9b.png">  

**data** : This layer deals with reading the data from files & creating the entities to be used by the application. It is easy to add new entities without caring about other layers, if there are more files.    

**service** : This layer provides a service to aggregate or combine the entities in a more meaningful way which can be used by the domain layer. Right now the application has only single service related to events but based on requirements more services can be added as necessary.  
_It should be noted that for the specified requirements, we could have merged this layer with data but I preferred to keep this separate for better extensible & maintainable design_   

**domain** : This layer contains the domain ideally resonating the terminology with the requirements. A domain can use single or multiple services based on the necessity. Currently there are 2 domains user & repo. If there comes more requirements, a new domain can be created in this layer.   

**command** : This layer contains all the commands & the respective handlers. They can make use of single or multiple domains to return the required output.  

- **Use of Heap**  
`Heap` data structure is used in the domain layer so as to return the top N elements. This is important for application's efficiency.  
`Less` function of the heap decides how to compare the elements. This function is being created in the command layer & domain layer use the given function to push or pop the entries from the heap.  
Creating a function dynamically in `command` layer (based on command `flags`) & passing it to `domain` layer, makes this application really flexible & extensible.      

- **Nested Sort Function**  
The application supports providing multiple `sort` fields for one of the command.  
The `sort` function (or `Less` function) is generated based on the these sort fields along with honoring the order in which the sort fields are provided.   
It is generated by starting from the last sort field & then wrapping it with function of sort field before it.      
This is done in a loop until the first sort field is reached.  
This may seem complicated at first but once you understand how it works, it feels trivial. Also, the flexibility it provides is very significant.  

- **Committing Data files to Repository**  
Committing data files to a repository is not a recommended approach.  
Ideally, a person cloning the repo should have data files on their machine.  
But I am committing the files in the repo, so that the person not having data files on their machine can also run this application. If you want to use the data files which are on your machine, you can provide the path of that directory.   

## Possible Enhancements
Due to time constraints, I could not work on following but it can be noted for future enhancements

- Structuring command layer into more logical separation
- Integration Tests   