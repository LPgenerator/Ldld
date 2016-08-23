## Client API interface

Api 

    export API='curl -X POST --user ldl:7eNQ4iWLgDw4Q6w -H "Accept: application/json" -s -i'
    export API_HOST="http://49.44.44.44:9090"


### Management

Get all available images

    $API $API_HOST/images


Pull image

    $API -d 'name=web' $API_HOST/pull


Show all commits for image

    $API -d 'name=web' $API_HOST/log


### CT

Create CT

    $API -d 'template=web&name=web-1' $API_HOST/create


Start CT

    $API -d 'name=web-1' $API_HOST/start


Stop CT

    $API -d 'name=web-1' $API_HOST/stop


Info about CT

    $API -d 'name=web-1' $API_HOST/info


Get all CT's

    $API $API_HOST/list


Destroy CT

    $API -d 'name=web-1' $API_HOST/destroy


Freeze CT

    $API -d 'name=web-1' $API_HOST/freeze


Unfreeze CT

    $API -d 'name=web-1' $API_HOST/unfreeze


Execute cmd at CT

    $API -d 'name=web-1&command=ls' $API_HOST/exec


### Settings

AutoStart after reboot

    $API -d 'name=web-1&value=1' $API_HOST/autostart


Port forwarding

    $API -d 'name=web-1&value=8080:80' $API_HOST/forward


Memory limit

    $API -d 'name=web-1&value=512' $API_HOST/memory


Swap limit

    $API -d 'name=web-1&value=1024' $API_HOST/swap


Static IP

    $API -d 'name=web-1&value=fix' $API_HOST/ip


CPU limits

    $API -d 'name=web-1&value=1' $API_HOST/cpu


Processes limits

    $API -d 'name=web-1&value=32' $API_HOST/processes


Network limits

    $API -d 'name=web-1&value=veth2P6NB7 5' $API_HOST/network
