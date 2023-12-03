#import "@preview/cetz:0.1.0"

#let half = 0.5

// note, anchor outbound values to x+1.1
#let and-gate((x, y), body) = {
    import cetz.draw: *

    let x_offset = 0.7
    let y_offset = 1
    

    line((x+x_offset, y), (x, y), (x, y+1), (x+x_offset, y+1), stroke: 0.7pt)
    content((x+half, y+half), body)
    arc((x+0.7, y), start: -30deg, stop: 30deg, radius: (2.7, 1), stroke: 0.7pt)
} 

// note, anchor outbound values to x+1.1
#let or-gate((x, y), body) = {
    import cetz.draw: *

    let x_offset = 0.7
    let y_offset = 1
    

    line((x+x_offset, y), (x, y), stroke: 0.7pt)
    line((x, y+1), (x+x_offset, y+1), stroke: 0.7pt)
    content((x+half, y+half), body)
    arc((x+x_offset, y), start: -30deg, stop: 30deg, radius: (2.7, 1), stroke: 0.7pt)
    arc((x, y), start: -30deg, stop: 30deg, radius: (2.7, 1), stroke: 0.7pt)
} 

// note, anchor outbound values to x+1.175, y+half
#let nand-gate((x, y), body) = {
    import cetz.draw: *

    let x_offset = 0.7
    let y_offset = 1
    

    line((x+x_offset, y), (x, y), (x, y+1), (x+0.7, y+1), stroke: 0.7pt)
    content((x+half, y+half), body)
    arc((x+0.7, y), start: -30deg, stop: 30deg, radius: (2.7, 1), stroke: 0.7pt)
    circle((x+1.175, y+half), radius:0.075, stroke: 0.7pt)
} 

// note, anchor outbound values to x+0.625
#let not-gate((x, y), body) = {
    import cetz.draw: *

    line((x, y+0.2), (x,y - 0.2), (x+half, y), (x, y+0.2), stroke: 0.7pt)
    circle((x+0.575, y), radius: 0.075, stroke: 0.7pt)
}

// note, inbound anchors at (x, y+0.5) and (x, y+1)
// outbound anchors at (x+0.5, y+0.75)
#let mux((x, y)) = {
    import cetz.draw: *

    line((x, y), (x+half, y+0.25), (x+half, y+1.25), (x, y+1.5), (x,y), stroke: 0.7pt)
    content((x+0.25, y+0.5), "1")
    content((x+0.25, y+1), "0")
}