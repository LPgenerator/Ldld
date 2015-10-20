## Server commands

    create     <ct-name> <template>   Creates a container.
    clone      <ct-src> <ct-dst>      Clone a new container from an existing one.
    commit     <ct-name>              Snapshot an existing container.
    push       <ct-name>              Push container diff to repository.
    export     <ct-name> <u@ssh.dst>  Export to remote client without repo sharing
    share                             Share local repository.    
    log        <ct-name>              List an existing snapshots for container.

    start      <ct-name>              Run container.
    stop       <ct-name>              Stop a container.
    freeze     <ct-name>              Freeze all the container's processes.
    unfreeze   <ct-name>              Thaw all the container's processes.
    destroy    <ct-name>              Destroy a container.
    info       <ct-name>              Query information about a container.
    attach     <ct-name>              Enter into a running container.
    exec       <ct-name> <cmd>        Execute the command inside CT.
    list                              List the containers existing on the system.
    help, h                           Shows a list of commands or help for one command
