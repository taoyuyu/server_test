#include <sys/types.h>
#include <sys/socket.h>
#include <stdio.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <fcntl.h>
#include <sys/shm.h>
#include <pthread.h>
#include <string.h>

#include <iostream>

#define MYPORT  5555
#define QUEUE   20
#define BUFFER_SIZE 1024

void get_load(char buf[]);
void system_call(char* cmd, char buf[]);
void* update_load(void* args);
char* get_target(char* p, char d[], int num);
int create_response(char response[], char msg[], int start);

pthread_mutex_t mutex;
static char load_of_cpu[8];
static char load_of_IO[8];
static char load_of_mem[8];

bool running = false;

int main()
{
    // new a thread to update load;
    memset(load_of_cpu,'\0',sizeof(load_of_cpu));
    memset(load_of_mem,'\0',sizeof(load_of_mem));
    pthread_t update_load_thread;
    running = true;
    int err = pthread_create(&update_load_thread, NULL, update_load, NULL);
    if (err != 0){
        printf("can't create thread: %s\n", strerror(err));
        exit(err);
    }

    ///定义sockfd
    int sockfd = socket(AF_INET,SOCK_STREAM, 0);

    ///定义sockaddr_in
    struct sockaddr_in server_sockaddr;
    server_sockaddr.sin_family = AF_INET;
    server_sockaddr.sin_port = htons(MYPORT);
    server_sockaddr.sin_addr.s_addr = htonl(INADDR_ANY);

    ///bind，成功返回0，出错返回-1
    if(bind(sockfd,(struct sockaddr *)&server_sockaddr,sizeof(server_sockaddr))==-1)
    {
        perror("bind");
        running = false;
        exit(1);
    }

    ///listen，成功返回0，出错返回-1
    if(listen(sockfd,QUEUE) == -1)
    {
        perror("listen");
        running = false;
        exit(1);
    }

    ///客户端套接字
    char buffer[2];
    struct sockaddr_in client_addr;
    socklen_t length = sizeof(client_addr);

    char response[24];
    while(true){
        ///成功返回非负描述字，出错返回-1
        int conn = accept(sockfd, (struct sockaddr*)&client_addr, &length);
        std::cout << inet_ntoa(client_addr.sin_addr) << std::endl;
        if(conn<0)
        {
            perror("connect");
            exit(1);
        }
        memset(buffer,'\0',sizeof(buffer));
        recv(conn, buffer, sizeof(buffer),0);

        memset(response, '\0', sizeof(response));
        int index = 0;
        // cpu + IO + mem;
        index = create_response(response, load_of_cpu, index);
        index = create_response(response, load_of_IO, index);
        index = create_response(response, load_of_mem, index);

        pthread_mutex_lock(&mutex);
        send(conn, response, sizeof(response), 0);
        pthread_mutex_unlock(&mutex);
        // close connection;
        close(conn);
    }
    close(sockfd);
    running = false;
    //wait update_load_thread end;
    pthread_join(update_load_thread, NULL);
    return 0;
}
int create_response(char response[], char msg[], int start){
    int index = start;
    for(int i = 0; i < sizeof(msg); i++){
        if(msg[i] == '\0'){
            response[index++] = '#';
            break;
        }
        response[index++] = msg[i];
    }
    return index;
}

char* get_target(char* p, char d[], int num){
    char* t_str = strtok(p,d);
    while(t_str){
        if(num == 0) return t_str;
        num --;
        t_str = strtok(NULL, d);
    }
}

void system_call(char* cmd, char buf[]){
    memset(buf, '\0', sizeof(buf));
    FILE* pf = popen(cmd, "r");
    fread(buf, BUFFER_SIZE, 1, pf);
    pclose(pf);
}

void* update_load(void* args){
    char buf[BUFFER_SIZE];
    while(running){
        system_call("iostat -x 1 2", buf);
        char d[] = "\n";
        int tag = 0;
        char d1[] = " *";
        char *target1, *target2;
        char* p = strtok(buf, d);
        while(p){
            if (tag == 6){
                target1 = p;
            }
            if (tag == 8){
                target2 = p;
                break;
            }
            tag++;
            p = strtok(NULL, d);
        }

        target1 = get_target(target1, d1, 0);
        target2 = get_target(target2, d1, 11);

        pthread_mutex_lock(&mutex);
        memcpy(load_of_cpu, target1, strlen(target1));
        memcpy(load_of_IO, target2, strlen(target2));
        pthread_mutex_unlock(&mutex);

        system_call("free -k | tail -2 | head -1|awk '{print $3}'", buf);
        double mem_used = atof(buf);

        system_call("free -k | tail -3 | head -1|awk '{print $2}'",buf);
        double mem_total = atof(buf);

        pthread_mutex_lock(&mutex);
        sprintf(load_of_mem,"%2.2f",mem_used/mem_total*100);
        pthread_mutex_unlock(&mutex);
    }
}
