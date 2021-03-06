from threading import Thread
import threading

i = 0
lock = threading.Lock();

def someThreadFunction1():
	lock.acquire()
	global i	
	for x in range(1000000):
		i+=1
	lock.release()

def someThreadFunction2():
	lock.acquire()
	global i	
	for x in range(999999):
		i-=1
	lock.release()

def main():
	someThread1 = Thread(target = someThreadFunction1, args = (),)
	someThread1.start()
	someThread2 = Thread(target = someThreadFunction2, args = (),)
	someThread2.start()

	someThread1.join()
	someThread2.join()
	print (i)

main()
