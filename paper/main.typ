#import "@preview/cetz:0.1.0"

#import "template.typ": *
#import "circuits.typ": *
#import "@preview/tablex:0.0.6": tablex, cellx, rowspanx, colspanx
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
Data and the calculations performed on data are physically separated in modern computing devices. Data is moved from a storage media to closer and closer locations before it is finally used for computation before being places back into memory. This paradigm of computing, the Von-Neumann paradigm, has been instrumental in the development of modern computing solutions. As processing has sped up though, this transfer of information from storage to processing has begun to form a bottleneck which limits the computational speeds which can be achieved by state-of-the-art devices. Specialized processing units, such as Graphics Processing Units and other purpose-built hardware, often attempts to side-step the problem by increasing the potential throughput of information; however, the theoretical problem of moving data around is not solved, only mitigated, by these solutions. An alternative paradigm to the von-Neumann architecture could be performing the computation at the same place the data is stored. This computing paradigm is aptly referred to as "in-memory computation."

Memristor Aided Logic (MAGIC) is an emerging computing paradigm making use of parallel, write-based systems to perform calculation in-memory @kvatinsky2014magic. This requires scheduling operations for the computation; however, the scheduling order, upon execution, may have substantially differing memory footprint requirements. State-of-the-art solutions model this dependence as a graph problem and perform scheduling as a graph covering problem. In this work, we characterize a number of properties of these evaluations graphs and apply those observations to the development of a genetic algorithm which produces reductions in the memory footprint of execution between 14% in the worst case and 26% in the best case when compared to standard algorithmic approaches

= Prior Works
/* one column */
// TODO
// - discuss baseline of in memory computing
// - lay out current SotA

Within the realm of the MAGIC, application of formal design principles has continued with advances in technology mapping @rashedautomated leading to advances in area, energy, and operation counts.

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

// TODO: Formally define dependent/dependance in this context

== Cost Metric <score-metric>

The memory footprint, in other words the cost, of a sequence $S$ derived from an adder DAG is defined as the number of memory cells that are need to evaluate the entire DAG. Determining the cost of a sequence can be done in linear time by simulating the execution sequence as a series of instructions. A memory bank is initialized with the labels of each vertex in the DAG with no parents (the input vertices). Then, the next element in the sequence is put into a free spae in the memory bank, expanding the width if necessary. This vertex is then marked as processed and all vertex labels which have all of their children marked as processed are removed from the memory bank. The maximum width at any point in this process is defined as: $"cost"(S)$.

= Genetic Algorithm

Genetic algorithms seek to emulate the processes of evolutionary biology observed in the physical world. They do this by modeling populations, natural selection, and reproduction over a series of "generations." An optimal solution, whether globally optimal or locally optimal, is found over time through selecting only the most viable samples from the population for inclusion in subsequent generations. 

Four common phases make up genetic algorithms: encoding, selection, crossover, and mutation @mirjalili2019genetic. This work uses a topological sorting of the DAG for the encoding. It uses a simple rank selection mechanism based on the score metric defined in @score-metric. For crossover, single point crossover based on analysis of the graph is employed, and swapping is used to mutate the sorting when mutation is applied. Each element of the algorithm is detailed in this section.

== Encoding and Selection

Encoding and selection are both informed by the underlying DAG, but make less explicit use of graph analysis than the other phases of the algorithm. The encoding of a sequence is a straightforward topological sort of the vertex labels in the graph. Trivially, a valid execution sequence can be attained using breadth first search. Sequences are therefore sorted according to the cost of the sequence for selection. During the selection phase, only the best performing sequences (those with minimal score) are selected to be present in the next evolutionary generation. 

== Crossover

An observation of the DAGs within this problem domain will quickly illustrate that there is a "cost-maximizing" step from which the cost of the remainder of execution is monotonically decreasing. This implies that for a block in the DAG which occurs after the cost-maximizing step occurs other values could remain in memory without increasing the cost of the sequence. Therefore, modifying the execution order of the block can be done without necessarily incurring a cost increase. However, the larger this block is, the more likely it is to include this cost-maximizing step (or to become this step by violating the dependence requirements of the DAG).

With these two ideas in mind, a single-point crossover mechanism is a natural fit to this problem. While it will occasionally break a block during the crossover, due to the nature of the DAG and this algorithm's place in the larger genetic algorithm this crossover mechanism could behave in a manner characterized as a search for an optimal order of executing blocks in the post-cost-maximizing phase of execution. As long as their order doesn't impact the overal cost of the sequence, the problem imposes no specific constraints on its performance. Therefore a crossover mechanism such as this will also fit naturally with the mutation mechanism employed for cost minimization.

== Mutation

A naive approach to mutation for this problem would be to randomly select two vertices in the sequence and swap their positions. This is a straightforward solution to implement, and it leads to sequence synthesis which performs similarly in regards to the sequence cost as the greedy algorithm presented in the literature and previously referred to in this work. However, in the context of the DAG modeling the problem, this strategy becomes more evidently sub-optimal. Two random vertices are likely to have an ancestor-descendant relationship of some sort, which implies that their swapping would create an invalid execution sequence. 

Instead of swapping vertices randomly in the execution sequence, we instead compute all of the peers of each vertex where a pair of vertices are peers if they have intersecting parents and children but are not directly dependant on one another. See @triangular-family for a visual depiction of this relationship. Extending this logic, a swapping of these two elements in the sequence can be seen as determining which of the blocks of the graph should be prepared for evaluation earlier in the sequence. 

When paired with the block-targeting crossover mechanism, this leads to an algorithm which is responsive to the underlying data modeled by the DAG. This leads to both faster convergence on a valid execution sequence than random mutation, but also improves cost in comparison with the greedy scheduling algorithm.

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
    caption: [The overlapping nature of children and parents illustrated visually.]
) <triangular-family>

= Evaluation
/* 1 column */
// TODO:
// - Discuss the software capabilities
// - Discuss the execution times

== Data Set
/* 0.5 columns */
// TODO: Describe the data and what it represents
The dataset evaluated in this work is synthesized from n-bit adder specifications using NOR and NAND gates. The netlist for the adder is modeled as a DAG and described using a proprietary format. This is parsed into a general DAG for evaluation. Six different adder widths were analyzed representing all powers of two less than 64 (i.e. [1, 2, 4, 8, 16, 32]). Execution sequences for each DAG were also provided. The memory footprints of these sequences are evaluated in the same way as candidate solutions will be evaluated in this work to determine a baseline memory footprint to compare future solutions against. @sota-footprint lists each adder's size and the best provided sequence's memory footprint. Complete BLIF specification adders come from the work of Rashed et al. @rashed2022logic.

#figure(
    box[
        #v(10pt)
        #tablex(
            columns: 2,
            inset: 4pt,
            align: center,
            [Adder Width], [Mem. Footprint],
            [1], [6],
            [2], [9],
            [4], [14],
            [8], [26],
            [16], [50],
            [32], [91]
        )
        #v(10pt)
    ],
    caption: [The memory footprints of adders with varying width discovered through the greedy algorithm @rashed2022logic.]
) <sota-footprint>

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