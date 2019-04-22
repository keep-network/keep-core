while [ 1 ];
  do echo -e "HTTP/1.1 200 OK\n\nenode://$(bootnode -writeaddress --nodekey=/data/boot.key)@$(POD_IP):30301" | nc -l -v -p 80 || break;
done;