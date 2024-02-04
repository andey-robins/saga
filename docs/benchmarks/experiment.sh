# iterature through each .blif file in the directory and run the command to evolution capture data
for file in ./docs/benchmarks/blif/*.blif; do
    ./magical -evolve -graph $file -out ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e500_mu20_s1.seq -pop 2000 -epsilon 500 -mutation 0.2 -verbose 2> ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e500_mu20_s1.out
done

for file in ./docs/benchmarks/blif/*.blif; do
    ./magical -evolve -graph $file -out ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e5000_mu20_s1.seq -pop 2000 -epsilon 5000 -mutation 0.2 -verbose 2> ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e5000_mu20_s1.out
done