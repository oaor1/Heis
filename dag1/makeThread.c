//gcc -std=gnu99 -Wall -g -o makeThread makeThread.c -lpthread

#include <pthread.h>
#include <stdio.h>

int i = 0;
void* someThreadFunction1(){
	for (int j=0; j <10000000; j++){
		i++;
	}
	return NULL;
}

void* someThreadFunction2(){
	for (int j=0; j <10000000; j++){
		i--;
	}
	return NULL;
}

int main(){
	pthread_t someThread1;
	pthread_t someThread2;
	pthread_create(&someThread1, NULL, someThreadFunction1, NULL);
	// Arguments to a thread would be passed here ---------^
	pthread_create(&someThread2, NULL, someThreadFunction2, NULL);

	pthread_join(someThread1, NULL);
	pthread_join(someThread2, NULL);
	printf(i);
	return 0;
}
