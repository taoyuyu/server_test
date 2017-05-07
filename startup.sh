killall server_test
killall a.out

go build server_test
nohup ./server_test start & ./get_load_server/a.out &
