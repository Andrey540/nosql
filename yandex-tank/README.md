docker run --rm -it -v $(pwd)/ammo:/ammo -v $(pwd)/logs:/logs --network="host" gtrafimenkov/yandex-tank -c /ammo/load.ini

docker run --rm -v $(pwd)/ammo:/var/loadtest -v $(pwd)/logs:/var/loadtest/logs --network="host" -it direvius/yandex-tank