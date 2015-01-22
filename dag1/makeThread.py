from threading import Thread

i = 0

def someThreadFunction1():
	global i	
	for x in range(1000000):
		i+=1

def someThreadFunction2():
	global i	
	for x in range(1000000):
		i-=1


def main():
	someThread1 = Thread(target = someThreadFunction1, args = (),)
	someThread1.start()
	someThread2 = Thread(target = someThreadFunction2, args = (),)
	someThread2.start()

	someThread1.join()
	someThread2.join()
	print (i)

main()
