#import "@preview/cetz:0.1.0"

#import "template.typ": *
#import "circuits.typ": *
#show: ieee.with(
  title: "MAGICAL - Genetic Algorithms for More Efficient In-memory Computation",
  abstract: [
    #lorem(150)
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
/* one column */
Data and the calculations performed on data are physically separated in modern computing devices. Data is moved from a storage media to closer and closer locations before it is finally used for computation before being places back into memory. This paradigm of computing, the Von-Neumann paradigm, has been instrumental in the development of modern computing solutions. As processing has sped up though, this transfer of information from storage to processing has begun to form a bottleneck which limits the computational speeds which can be achieved by state-of-the-art devices. Specialized processing units, such as Graphics Processing Units and other purpose-built hardware, often attempts to side-step the problem by increasing the potential throughput of information; however, the theoretical problem of moving data around is not solved, only mitigated, by these solutions. An alternative paradigm to the Von-Neumann architecture could be performing the computation at the same place the data is stored. This computing paradigm is aptly referred to as "in-memory computation."

Memristor Aided Logic (MAGIC) is an emerging computing paradigm making use of parallel, write-based systems to perform calculation in-memory @kvatinsky2014magic. This requires scheduling operations for the computation; however, the scheduling order, upon execution, may have substantially differing memory footprint requirements. State-of-the-art solutions model this dependence as a graph problem and perform scheduling as a graph covering problem. In this work, we characterize a number of properties of these evaluations graphs and apply those observations to the development of a genetic algorithm which produces reductions in the memory footprint of execution between 14% in the worst case and 26% in the best case when compared to standard algorithmic approaches

= Prior Works
/* one column */
// TODO
// - discuss baseline of in memory computing
// - lay out current SotA


= Problem Specification

/* one column */
// TODO:
// - Formally lay this out as a graph problem
// - Highlight the mapping from engineering domain to graph theory

In-memory computation can be modeled as an execution graph. Translated from traditional combinational logic, each vertex in the graph corresponds with a boolean logic gate. Within the memristor array, the input values to the logic gate can be selected by the write signal before being written to an empty location in the memristor array. An edge exists in the graph from a vertex to another if they have this depednece relation. As an example, for evaluating the boolean function $f = a b' + c$, @mapping-spaces details the transformation from function to circuit to graph to in-memory computation. Beginning with a boolean function $f$, existing processes are highly capable of mapping this to a minimal circuit. This circuit can then be transformed into the graphs under discussion in this work by assigning each gate a vertex and making each wire an edge in the graph. The final step illustrated in @mapping-spaces is the evaluation using in-memory computation and illustrates the requirement of more memory cells than inputs to the function. The objective of this work is to minimize the extra memory needed.

#figure(
    box[
        #cetz.canvas({
            import cetz.draw: *
            let r = 0.4
            content((1.2,1.2), stroke: red)[CHANGE ME]
            circle((1.5, 3.5), name: "1", radius: r, fill: yellow)
            content((1.5, 3.5))[$m_1$]
            circle((3.5, 2.5), name: "2", radius: r)
            content((3.5, 2.5))[$m_2$]
            circle((4, 1), name: "3", radius: r)
            content((4, 1))[$m_3$]
            circle((2.5, -0.5), name: "4", radius: r)
            content((2.5, -0.5))[$m_4$]
            circle((0.5, -0.5), name: "5", radius: r, fill: yellow) 
            content((0.5, -0.5))[$m_5$]
            circle((-1, 1), name: "6", radius: r, fill: yellow)
            content((-1, 1))[$m_6$]
            circle((-0.5, 2.5), name: "7", radius: r)
            content((-0.5, 2.5))[$m_7$]

            line("1.bottom", "7.right")
            line("1.bottom", "2.left")
            line("1.bottom", "4.top")
            line("2.left", "7.right")
            line("2.left", "6.right")
            line("2.left", "5.top")
            line("2.left", "4.top")
            line("2.left", "3.left")
            line("3.left", "7.right")
            line("3.left", "6.right")
            line("3.left", "5.top")
            line("4.top", "6.right")
            line("4.top", "5.top")
            line("5.top", "7.right")  
        })
        #v(20pt)
    ]
    ,
    caption: [The MIS graph constructed from the covering table presented in the assignment.]
) <mapping-spaces>

Formally, the interdependence between computational nodes in the execution graph is a directed acyclic graph (DAG) in which the children of a vertex must not be executed until that vertex is executed. Thus, for any vertex, it can be viewed as both the root of 

== Data Set
/* 0.5 columns */
// TODO: Describe the data and what it represents

= Genetic Algorithm
/* 1.5 columns */
// TODO:
// - Outline specifics of this algorithm
// - Describe triangular mapping properties

= Evaluation
/* 1 column */
// TODO:
// - Discuss the software capabilities
// - Discuss the execution times

= Results
/* 1 column */
// TODO:
// - Compare memory footprints to SotA
// - Potentially discuss execution results

= Discussion
/* 0.5 column */
// Highlight how modeling this as a graph analysis allowed for these results

= Conclusion
/* 0.5 column */
// Future work, final takeaways, restate, etc.