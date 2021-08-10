#!/bin/bash
#Backup by namespace

usage()
{
    echo "Usage: sh $0 namespace"
}

getmark() {
    echo "***************************************""$1""***************************************"
}

#Judge cmd
if [ $# -ne 1 ]
then
        usage
        exit 1
fi

namespace=$1

#Judge if namespace exists
getmark "Judge if namespace $namespace exists"
kubectl get namespace $namespace
if [ $? -ne 0 ]
then
   echo "命名空间: $namespace 不存在,退出"
   exit 1
fi

backupcmd=/home/opts/velero-v1.5.3-linux-amd64/velero

#删除旧的备份
olddatedir=`date +%Y%m%d -d '7 day ago'`
oldbackupname="$namespace"-backup-"$olddatedir"
getmark "Delete old backup, namespace: $namespace,oldbackupname: $oldbackupname"
$backupcmd backup delete $oldbackupname --confirm

#执行备份
datedir=`date +%Y%m%d`
backupname="$namespace"-backup-"$datedir"
getmark "Backup namespace: $namespace,backupname: $backupname"
$backupcmd backup get $backupname
if [ $? -ne 0 ]
then
   $backupcmd backup create $backupname --include-namespaces $namespace
else
   echo "$backupname exists already."
   #$backupcmd backup delete $backupname --confirm
   #$backupcmd backup create $backupname --include-namespaces $namespace
fi
