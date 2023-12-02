#import "template.typ": *
#show: ieee.with(
  title: "MAGICAL - Genetic Algorithms for More Efficient In-memory Computation",
  abstract: [

  ],
  authors: (
    (
      name: "Andey Robins",
      department: [Dept. of Electrical and Computer Engineering],
      organization: [University of Central Florida],
      location: [Orlando, USA],
      email: "andey.robins@ucf.edu"
    ),
  ),
  index-terms: ("in-memory computation", "artificial intelligence", "genetic algorithms", "computer aided design"),
  bibliography-file: "refs.bib",
)

= Introduction
Memristor Aided Logic (MAGIC) is an emerging computing paradigm making use of parallel, write-based systems to perform calculation in-memory @kvatinsky2014magic. This requires scheduling operations for the computation; however, the scheduling order, upon execution, may have substantially differing memory footprint requirements. State-of-the-art solutions model this dependence as a graph problem and perform scheduling as a graph covering problem. In this work, we characterize a number of properties of these evaluations graphs and apply thos observations to the development of a genetic algorithm which produces reductions in the memory footprint of execution between 14% in the worst case and 26% in the best.

= Prior Works

= Genetic Algorithm

= Evaluation

= Results

= Discussion

= Conclusion