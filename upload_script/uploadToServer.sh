#!/bin/bash
# Скрипт загружает ./dist папку на тестовый сервер
. .env.local

if [[ ! $SERVER_USER_NAME ]]
then
    echo "Не указан SERVER_USER_NAME в .env "
    exit
fi
if [[ ! $SERVER_PUB_KEY_PATH ]]
then
    echo "Не указан SERVER_PUB_KEY_PATH в .env "
    exit
fi
if [[ ! $SERVER_APP_NAME ]]
then
    echo "Не указан SERVER_APP_NAME в .env "
    exit
fi

_isYes=false
while getopts y flag
do
    case "${flag}" in
        y) _isYes=true;;
    esac
done

echo 'Параметры из .env.local'
echo 'SERVER_USER_NAME='$SERVER_USER_NAME
echo 'SERVER_PUB_KEY_PATH='$SERVER_PUB_KEY_PATH
echo 'SERVER_APP_NAM='$SERVER_APP_NAME
echo ' '

_appPath=/var/www/projects/app/$SERVER_APP_NAME

echo "Файлы на сервере:"
echo "======="
ssh $SERVER_USER_NAME@88.99.160.4 -p 58533 'cd '$_appPath'/ && pwd && ls'
if [[ $_isYes == false ]]
then
    read -p "Удалить их и загрузить новые (y/n)? " choice
    case "$choice" in 
    y|Y ) echo "Аминь!";;
    n|N ) echo "Возвращайтесь когда будете готовы"; exit;;
    * ) echo "y или n, написано же"; exit;;
    esac
fi

echo "Удаляем файлы"
ssh $SERVER_USER_NAME@88.99.160.4 -p 58533 'rm -rf '$_appPath'/* '$_appPath'/.htaccess'

echo "Архивируем dist папку"
cd ./dist
chmod -R 775 *
rm -f dist.zip
zip -r dist.zip .

echo "Теперь копируем архив на сервер"
scp -i $SERVER_PUB_KEY_PATH -P 58533 dist.zip $SERVER_USER_NAME@88.99.160.4:$_appPath
# scp -i ~/.ssh/motiv_rsa.pub -P 58533 -r ./dist maxi@88.99.160.4:/var/www/projects/app/'$SERVER_APP_NAME'/test

echo "Распакуем его и удаляем архив"
ssh $SERVER_USER_NAME@88.99.160.4 -p 58533 'bash -c "cd '$_appPath'/ && unzip dist.zip && rm -f dist.zip"'

echo "Поздравляю! Идём и проверяем что всё открывается!"
echo "https://"$SERVER_APP_NAME".app.motivationportal.ru"

cd ../
