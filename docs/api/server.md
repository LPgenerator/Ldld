## Server API interface


Api 

    export API='curl -X POST --user ldl:7eNQ4iWLgDw4Q6w -H "Accept: application/json" -s -i'
    export API_HOST="http://127.0.0.1:9090"


### Management

Image log

    $API -d 'name=web' $API_HOST/log


Commit image

    $API -d 'name=web' $API_HOST/commit


Push image

    $API -d 'name=web' $API_HOST/push



### CT

Create CT

    $API -d 'name=web' $API_HOST/create


Start CT

    $API -d 'name=web' $API_HOST/start


Stop CT

    $API -d 'name=web' $API_HOST/stop


Info about CT

    $API -d 'name=web' $API_HOST/info


Get all CT's

    $API $API_HOST/list


Clone all CT

    $API -d 'from=web&to=cron' $API_HOST/clone


Destroy CT

    $API -d 'name=web' $API_HOST/destroy


Freeze CT

    $API -d 'name=web' $API_HOST/freeze


Unfreeze CT

    $API -d 'name=web' $API_HOST/unfreeze


Execute cmd at CT

    $API -d 'name=web&command=ls' $API_HOST/exec
