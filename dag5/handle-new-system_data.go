sjekkliste


gjennoppstaa 
	spør etter system_data

motta system_data //(etter dødsfall)
	merge med eventuell egen versjon av system_data

init
	spør etter system_data


Vinn bud 
	send system_data_update


håndter ordre 
	send system_data_update


motta bestilling 
	send system_data_update

motta system_data_update
	legg til i egen system_data



aution handler: 
	
	ved ny bestilling sender lokal heis ut en system_data_update der 
	deadline blir presisert og skal legges til i aution_q up/down

	while up_q != 0:
		calculate_local_bid
		bordcast local bid
		

