## Ldld

`Ldld` is a simple and open source tool for shipping and running distributed containers based on lxc. 
Work with your containers such as git from command line and with HTTP REST Api.


[![Build Status](http://ci.lpgenerator.ru/projects/7/status.png?ref=master)](http://ci.lpgenerator.ru/projects/7?ref=master)


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
    ldld pull web                               # download web image
    ldld log web                                # show commits for web
    ldld create web-1 web                       # create new CT from latest web version
    ldld create web-2 web                       # create new CT from latest web version
    ldld create web-3 web:0                     # create CT from first commit on web 
    ldld start web-1                            # start web-1
    ldld start web-3                            # start web-3
    ldld list                                   # show all CT's 
    ldld exec web-1 'ls hello.txt'              # check file exists
    ldld exec web-3 'ls hello.txt'              # file is not found on first commit
    

### Contributing

The official repository for this project is on [Github.com](https://github.com/LPgenerator/Ldld).

* [Development](docs/development/README.md)
* [Issues](https://github.com/LPgenerator/Ldld/issues)
* [Pull Requests](https://github.com/LPgenerator/Ldld/pulls)


### Requirements

This project is designed for the Linux operating system.

* Ubuntu >= 14.04
* LXC >= 1.0.7
* ZFS >= 0.6.5 || BtrFS >= 3.13 || OverlayFS >= 3.13

Originally designed for ZFS and BtrFS. OverlayFS is experimental. 
Lvm2 is not completed and wait help from community.


### Installation

* [Install on Linux](docs/installation/README.md)
* [Install development environment](docs/development/README.md)


### Support CT templates

We are support distros based on debian 

* Debian
* Ubuntu

_Note: For OverlayFS we are support only Ubuntu at this moment_


### License

GPLv3
