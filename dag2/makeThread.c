//gcc -std=gnu99 -Wall -g -o makeThread makeThread.c -lpthread

#include <pthread.h>
#include <stdio.h>
#include <semaphore.h>
#include <time.h>

clock_t begin, end;
double time_spent;


int i = 0;
sem_t iSem;

void* someThreadFunction1(){
	sem_wait(&iSem);
	for (int j=0; j <10000000; j++){
		i++;
	}
	sem_post(&iSem);
	return NULL;
}

void* someThreadFunction2(){
	sem_wait(&iSem);	
	for (int k=0; k <9999999; k++){
		i--;
	}
	sem_post(&iSem);
	return NULL;
}

int main(){

	begin = clock();
	
	

	sem_init(&iSem, 0, 1);
	
	pthread_t someThread1;
	pthread_t someThread2;
	pthread_create(&someThread1, NULL, someThreadFunction1, NULL);
	// Arguments to a thread would be passed here ---------^
	pthread_create(&someThread2, NULL, someThreadFunction2, NULL);

	pthread_join(someThread1, NULL);
	pthread_join(someThread2, NULL);
	
	end = clock();
	time_spent = (double)(end - begin) / CLOCKS_PER_SEC;
	printf("\n\n %i \n time spent = %f\n",i ,time_spent);
	return 0;
}
