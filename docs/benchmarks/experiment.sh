go build

echo "Running for epsilon 50"
for file in ./docs/benchmarks/blif/*.blif; do
    echo $file
    ./magical -evolve -graph $file -out ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e50_mu20_s2.seq -seed 2 -pop 2000 -epsilon 50 -mutation 0.2 -verbose 2> ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e50_mu20_s1.out
done

echo "Running for epsilon 500"
for file in ./docs/benchmarks/blif/*.blif; do
    echo $file
    ./magical -evolve -graph $file -out ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e500_mu20_s2.seq -seed 2 -pop 2000 -epsilon 500 -mutation 0.2 -verbose 2> ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e500_mu20_s1.out
done

echo "Running for epsilon 5000"
for file in ./docs/benchmarks/blif/*.blif; do
    echo $file
    ./magical -evolve -graph $file -out ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e5000_mu20_s2.seq -seed 2 -pop 2000 -epsilon 5000 -mutation 0.2 -verbose 2> ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e5000_mu20_s1.out
done