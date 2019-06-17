#!/bin/bash

SERVER=basename "$(PWD)"

PID=$(ps -e | grep ${SERVER} | awk '{print $1}') # 获取进程 PID


function waiting()
{
        i=0
        while [[ ${i} -le 10 ]]
        do
                for j in '\\' '|' '/' '-'
                    do
                        printf "%c%c%c%c%c ······ %c%c%c%c%c\r" \
                        "$j" "$j" "$j" "$j" "$j" "$j" "$j" "$j" "$j" "$j"
                        sleep 0.1
                    done
                        let i=i+4
        done
                printf "\n"
}

function start()
{
        echo ${PID}

#        if [[ "`pgrep ${SERVER} -u ${UID}`" != "" ]];then
#		        echo "${SERVER} already running"
#		        exit 1
#	    fi

        nohup ${BASE_DIR}/${SERVER}  server &>/dev/null &

        waiting

        PID=$(ps -e | grep ${SERVER} | awk '{print $1}') # 更新 PID
        echo "update pid:${PID}"

        # check status
	    if [[ "`pgrep ${SERVER} -u ${UID}`" == "" ]];then
		        echo "${SERVER} start failed"
		        exit 1
	    fi
}

function status()
{
        if [[ "`pgrep ${SERVER} -u ${UID}`" != "" ]];then
		        echo ${SERVER} is running
	    else
		        echo ${SERVER} is not running
	    fi
}

function stop()
{
        echo ${PID}
        kill -9 ${PID}

        waiting

        PID=$(ps -e | grep ${SERVER} | awk '{print $1}') # 更新 PID
        echo "update pid:${PID}"

        if [[ "`pgrep ${SERVER} -u ${UID}`" != "" ]];then
		        echo "${SERVER} stop failed"
		        exit 1
	    fi
}

case "$1" in
        'start')
        start
        ;;
        'stop')
        stop
        ;;
        'status')
        status
        ;;
        'restart')
        stop && start
        ;;
        *)
        echo "usage: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac
