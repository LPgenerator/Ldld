## Ldld

`Ldld` is a simple and open source tool for shipping and running distributed containers based on LXC. 
Work with your containers such as git from command line and with HTTP REST Api.


### Demo

    # Clone and up VM's
    git clone --depth 1 https://github.com/LPgenerator/Ldld.git; cd Ldld
    vagrant up --provider virtualbox
    
    # Master
    vagrant ssh master
    sudo -i
    ldld create web                             # create new basic CT (ubuntu by default)
    ldld start web                              # start CT
    ldld info web                               # check CT state
    ldld commit web                             # initial commit
    ldld exec web 'touch hello.txt'             # create new file
    ldld commit web                             # new commit with file
    ldld log web                                # show commits
    ldld push web                               # push all commits to repo
    
    # Slave
    vagrant ssh slave
    sudo -i
    ldld images                                 # show all available images on repo
    ldld pull web                               # download web images
    ldld log web                                # show all commits for web
    ldld create web-1 web                       # create new CT-1 instance from latest web version
    ldld create web-2 web                       # create new CT-2 instance from latest web version
    ldld create web-3 web:0                     # create new CT-2 instance from first commit 
    ldld start web-1                            # start web-1
    ldld start web-3                            # start web-3
    ldld list                                   # show all CT's 
    ldld exec web-1 'ls hello.txt'              # check file is exists on web-1 on latest version
    ldld exec web-3 'ls hello.txt'              # file not found on first commit
    

### Contributing

The official repository for this project is on [Github.com](https://github.com/LPgenerator/Ldld).

* [Development](docs/development/README.md)
* [Issues](https://github.com/LPgenerator/Ldld/issues)
* [Pull Requests](https://github.com/LPgenerator/Ldld/pulls)


### Requirements

This project is designed for the Linux operating system.

* Ubuntu >= 14.04
* LXC >= 1.0.7
* ZFS >= 0.6.5 || BtrFS >= 3.13


### Installation

* [Install on Linux](docs/installation/README.md)
* [Install development environment](docs/development/README.md)


### Support CT templates

We are support distros based on Debian 

* Debian
* Ubuntu



### License

GPLv3
