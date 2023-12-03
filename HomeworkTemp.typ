#let HomeworkHeader(class, prof, title, rest) = {
    box(height: 60pt)[
        #box(width: 1fr)[
            #text(size: 11pt, class) \
            #text(size: 10pt)[University of Central Florida] \
            #text(size: 10pt, smallcaps(prof)) \
            #line(length: 130pt, stroke: 0.7pt)
        ]
        #box(width: 1fr)[
            #align(center)[
                #text(size: 18pt, title) \
                #text(size: 11pt, smallcaps("Andey Robins")) \
                #v(15pt)
            ]
        ]
        #box(width: 1fr)[
            #align(right)[
                #text(size: 11pt)[
                    #datetime.today().display("[month repr:long] [day], [year]")
                ] \ \ \
                #line(length: 130pt, stroke: 0.7pt)
            ]
        ]
    ]
    rest
}