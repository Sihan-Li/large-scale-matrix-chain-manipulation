mkdir time
Mkdir output
go build -o executable proj3/calculate

for P in 5 20 50 100
do
    for M in 20 50 100 200
    do
    	for N in 1 4 8 12 16
    	do        
		TIMEFMT=$'\nreal\t%E\nuser\t%U\nsys\t%S'
	        echo p_timing_file"$P"_"$M"_"$N"
	        (time ./executable "$P" "$M" "$N") > output/out"$P"_"$M"_"$N".txt 2> time/p_timing_"$P"_"$M"_"$N"
        done
    done
done