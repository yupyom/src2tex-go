/* {\hrulefill} *


{\ % beginning of TeX mode

\centerline{\bf Towers of Hanoi (Rust)}

\begin{quote}

This program gives an answer to the following famous problem (towers of
Hanoi).
There is a legend that when one of the temples in Hanoi was constructed,
three poles were erected and a tower consisting of 64 golden discs was
arranged on one pole, their sizes decreasing regularly from bottom to top.
The monks were to move the tower of discs to the opposite pole, moving
only one at a time, and never putting any size disc above a smaller one.
The job was to be done in the minimum numbers of moves. What strategy for
moving discs will accomplish this optimum transfer?
\end{quote}

% end of TeX mode }

* {\hrulefill} */


/* {\hrulefill\ hanoi.rs \ \hrulefill} */


const ARRAY: usize = 8;		/* {\ disc の数 \hfill} */

static mut DISC: [[i32; 8]; 3] = [[0; 8]; 3];
					/* {\ disc に関するデータの置き場所
					   \hfill} */

fn init_array() {			/* {\ disc に関するデータの初期化
					   \hfill} */
    unsafe {
        for j in 0..ARRAY {
            DISC[0][j] = (ARRAY - j) as i32;
            DISC[1][j] = 0;
            DISC[2][j] = 0;
        }
    }
}

static mut COUNTER: i32 = 0;		/* {\ 移動回数カウンタ \hfill} */

fn print_result() {			/* {\ 結果の表示 \hfill} */
    unsafe {
        COUNTER += 1;
        println!("---#{}---", COUNTER);
        for i in 0..3 {
            print!("[{}] ", i);
            for j in 0..ARRAY {
                if DISC[i][j] != 0 {
                    print!("{} ", DISC[i][j]);
                } else {
                    break;
                }
            }
            println!();
        }
    }
}

static mut PTR: [usize; 3] = [0, 0, 0];
					/* {\ disc 移動用ポインタ（インデックス）
					   \hfill} */

fn move_one_disc(i: usize, j: usize) {	/* {\ 1枚の disc を pole $i$ から
					   pole $j$ に移動する \hfill} */
    unsafe {
        PTR[i] -= 1;
        DISC[j][PTR[j]] = DISC[i][PTR[i]];
        PTR[j] += 1;
        DISC[i][PTR[i]] = 0;
    }
}

fn move_discs(n: i32, i: usize, j: usize, k: usize) {
					/* {\ 上から $n$ 枚目までの disc
					   を、pole $i$ から pole $j$ に
					   \hfill} */
					/* {\ pole $k$ を経由して、移動する
					   \hfill} */
    if n >= 1 {
        move_discs(n - 1, i, k, j);	/* {\ 関数 {\tt move\_discs()}
					   の中で、さらに自分自身 \hfill} */
        move_one_disc(i, j);		/* {\ {\tt move\_discs()} が使われ
					   ている。このような \hfill} */
        print_result();			/* {\ 手法は、「再帰的呼びだし」
					   といわれる。 \hfill} */
        move_discs(n - 1, k, j, i);
    }
}

/* {\par\begin{center}

\includegraphics[scale=0.3]{hanoi1}\quad
\includegraphics[scale=0.3]{hanoi2}\end{center}


たとえば、関数 {\tt move\_discs(4, 0, 1, 2)} を呼び出すことは、
上図のような操作をすることに対応する。\hfill} */


fn main() {
    unsafe {
        PTR[0] = ARRAY;
        PTR[1] = 0;
        PTR[2] = 0;
    }

    init_array();
    move_discs(ARRAY as i32, 0, 1, 2);	/* {\ {\tt ARRAY} 枚の disc を
					   pole 0 から pole 1 に pole 2
					   \hfill} */
					/* {\ を経由して、移動する \hfill} */
}
