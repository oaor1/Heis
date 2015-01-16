from threading import Thread

i = 0

def someThreadFunction1():
	global i	
	while i<1000000:
		i+=1
	print ("first\n")

def someThreadFunction2():
	global i	
	while i<1000000:
		i+=1
	print ("second\n")


def main():
	someThread1 = Thread(target = someThreadFunction1, args = (),)
	someThread1.start()
	someThread2 = Thread(target = someThreadFunction2, args = (),)
	someThread2.start()

	someThread1.join()
	someThread2.join()
	print("ferdig!\n")

main()
