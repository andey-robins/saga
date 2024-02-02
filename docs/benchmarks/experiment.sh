# iterature through each .blif file in the directory and run the command to evolution capture data
for file in ./docs/benchmarks/blif/*.blif; do
    ./magical -blif -graph $file -out ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e1000_mu20_s1.seq -pop 2000 -epsilon 1000 -mutation 0.2 -verbose 2> ./docs/benchmarks/sequences/$(basename $file .blif)_p2000_e1000_mu20_s1.out
done