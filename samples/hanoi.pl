#!/usr/bin/perl
# {\hrulefill }
#
# {\ % beginning of TeX mode }
#
# {\ \centerline{\bf Towers of Hanoi (Perl)} }
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

# {\hrulefill\ hanoi.pl \ \hrulefill}

use strict;
use warnings;

my $ARRAY = 8;                          # disc count
my @disc;                               # disc data storage
for my $i (0..2) {
    $disc[$i] = [(0) x $ARRAY];
}

sub init_array {                        # initialize disc data
    for my $j (0..$ARRAY-1) {
        $disc[0][$j] = $ARRAY - $j;
        $disc[1][$j] = 0;
        $disc[2][$j] = 0;
    }
}

my $counter = 0;                        # move counter
my @ptr = (0, 0, 0);                    # stack pointers

sub print_result {                      # print current state
    $counter++;
    print "---${counter}---\n";
    for my $i (0..2) {
        print "[$i] ";
        for my $j (0..$ARRAY-1) {
            if ($disc[$i][$j] != 0) {
                print "$disc[$i][$j] ";
            } else {
                last;
            }
        }
        print "\n";
    }
}

sub move_one_disc {                     # move 1 disc: pole i to pole j
    my ($i, $j) = @_;
    $ptr[$i]--;
    $disc[$j][$ptr[$j]] = $disc[$i][$ptr[$i]];
    $ptr[$j]++;
    $disc[$i][$ptr[$i]] = 0;
}

sub move_discs {                        # move n discs
    my ($n, $i, $j, $k) = @_;          # pole i to pole j
    if ($n >= 1) {                      # via pole k
        move_discs($n - 1, $i, $k, $j);
        move_one_disc($i, $j);          # move_discs()
        print_result();                 # print result
        move_discs($n - 1, $k, $j, $i);
    }
}

init_array();
$ptr[0] = $ARRAY;
print_result();
move_discs($ARRAY, 0, 2, 1);
