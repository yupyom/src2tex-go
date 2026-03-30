#!/bin/sh
# {\hrulefill }
#
# {\ % beginning of TeX mode }
#
# {\ \centerline{\bf Towers of Hanoi (Shell)} }
#
# {\ \begin{quote} }
# {\ \noindent }
# {\ This program gives an answer to the following famous problem (towers of }
# {\ Hanoi). }
# {\ There is a legend that when one of the temples in Hanoi was constructed, }
# {\ three poles were erected and a tower consisting of 64 golden discs was }
# {\ arranged on one pole, their sizes decreasing regularly from bottom to top. }
# {\ The monks were to move the tower of discs to the opposite pole, moving }
# {\ only one at a time, and never putting any size disc above a smaller one. }
# {\ The job was to be done in the minimum numbers of moves. What strategy for }
# {\ moving discs will accomplish this optimum transfer? }
# {\ \end{quote} }
#
# {\ % end of TeX mode }
#
# {\hrulefill }

# {\hrulefill\ hanoi.sh \ \hrulefill}


ARRAY=8         # disc count
counter=0       # move counter

# initialize disc data
init_array() {
    for j in $(seq 0 $((ARRAY - 1))); do
        eval "disc_0_$j=$((ARRAY - j))"
        eval "disc_1_$j=0"
        eval "disc_2_$j=0"
    done
    ptr_0=$ARRAY
    ptr_1=0
    ptr_2=0
}

# print current state
print_result() {
    counter=$((counter + 1))
    echo "---${counter}---"
    for i in 0 1 2; do
        printf "[%d] " "$i"
        for j in $(seq 0 $((ARRAY - 1))); do
            eval "val=\$disc_${i}_${j}"
            if [ "$val" -ne 0 ]; then
                printf "%d " "$val"
            else
                break
            fi
        done
        echo
    done
}

# move 1 disc: pole i to pole j
move_one_disc() {
    i=$1; j=$2
    eval "ptr_i=\$ptr_$i"
    ptr_i=$((ptr_i - 1))
    eval "ptr_$i=$ptr_i"
    eval "ptr_j=\$ptr_$j"
    eval "disc_${j}_${ptr_j}=\$disc_${i}_${ptr_i}"
    ptr_j=$((ptr_j + 1))
    eval "ptr_$j=$ptr_j"
    eval "disc_${i}_${ptr_i}=0"
}

# move n discs
move_discs() {
    n=$1; i=$2; j=$3; k=$4                      # pole i to pole j
    if [ "$n" -ge 1 ]; then                     # via pole k
        move_discs $((n - 1)) "$i" "$k" "$j"
        move_one_disc "$i" "$j"                 # move_discs()
        print_result                            # print result
        move_discs $((n - 1)) "$k" "$j" "$i"
    fi
}

init_array
print_result
move_discs $ARRAY 0 2 1
