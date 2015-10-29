## Client commands

    images                                 List the images existing on the system.
    pull       <image>                     Pull and apply images updates from repository.
    import     <image>                     Import images and apply when used export command on master.
    create     <ct-name> <image>           Creates a container from specified image.
    migrate    <ct-name> <u@ssh.dst>       Migrate CT to a new host (include mounted folders and images).
    log        <ct-name>                   List an existing snapshots for container.
    run                                    Run Ldl http API.
    
    start      <ct-name>                   Run container.
    stop       <ct-name>                   Stop a container.
    freeze     <ct-name>                   Freeze all the container's processes.
    unfreeze   <ct-name>                   Thaw all the container's processes.
    destroy    <ct-name>                   Destroy a container.
    info       <ct-name>                   Query information about a container.
    attach     <ct-name>                   Enter into a running container.
    exec       <ct-name> <cmd>             Execute the command inside CT.
    list                                   List the containers existing on the system.
    
    autostart  <ct-name> <value>           Autostart after a reboot (0 or 1).
    forward    <ct-name> <src:dst>         Port forwarding (ip will be fixed).
    mount      <ct-name> <dir-name> <dst>  Create and mount new zfs dataset to CT.
    umount     <ct-name> <dir-name>        Unmount zfs dataset from CT (dataset is not deleted).
    cgroup     <ct-name> <key> <val>       NO DESCRIPTION
    memory     <ct-name> <value>           Memory limits. Set a maximum RAM (on MB).
    swap       <ct-name> <value>           Swap limits. Set a maximum swap (on MB).
    ip         <ct-name> <value>           Static IP. If value = 'fix', current IP will be fixed.
    cpu        <ct-name> <value>           CPU limits based on CPU shares.
    help, h                                Shows a list of commands or help for one command.
