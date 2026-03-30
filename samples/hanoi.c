/* {\hrulefill} *


{\ % beginning of TeX mode

\centerline{\bf Towers of Hanoi}

\begin{quote}
\noindent
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


/* {\hrulefill\ hanoi.c\ \hrulefill} */


#include <stdio.h>
#define ARRAY 8				/* {\ disc の数 \hfill} */

int disc[3][ARRAY];			/* {\ disc に関するデータの置き場所
					   \hfill} */

void init_array(void)			/* {\ disc に関するデータの初期化
					   \hfill} */
{
    int j;

    for (j = 0; j < ARRAY; ++j)		/* {\ [問題] この {\bf for} ループ
					   を抜け出した後、２次元配列
					   \hfill} */
					/* {\null{\tt disc[][]} にはどのよう
					   なデータが代入されることに
					   \hfill} */
					/* {\ なるか、具体的に述べよ。
					   \hfill} */
      {
	  disc[0][j] = ARRAY - j;
	  disc[1][j] = 0;
	  disc[2][j] = 0;
      }
}

void print_result(void)			/* {\ 結果の表示 \hfill} */
{
    static long counter = 0;		/* {\ [問題] ここで {\tt static}
					   と宣言したのは、なぜなのか
					   \hfill} */
					/* {\ 考えよ。もしも {\tt static}
					   を削除したら、何が起こる \hfill} */
					/* {\ か、実験をしてみよう。
					   \hfill} */
    int i, j;

    printf("---#%ld---\n", ++counter);
    for (i = 0; i <= 2; ++i)
      {
	  printf("[%d] ", i);
	  for (j = 0; j < ARRAY; ++j)
	    {
		if (disc[i][j] != 0)
		    printf("%d ", disc[i][j]);
		else
		    break;
	    }
	  printf("\n");
      }
}

static int *ptr[3];			/* {\ disc 移動用ポインタ \hfill} */

void move_one_disc(int i, int j)	/* {\ 1枚の disc を pole $i$ から
					   pole $j$ に移動する \hfill} */
{
    (*ptr[j]++) = (*--ptr[i]);		/* {\ [問題] {\tt ++} はポインタ
					   の後にあり、{\tt --} はポインタ
					   \hfill} */
					/* {\ の前にある。なぜ、そのように
					   しなければならない \hfill} */
					/* {\ のか説明せよ。\hfill} */
    *ptr[i] = 0;
}

void move_discs(int n, int i, int j, int k)
					/* {\ 上から $n$ 枚目までの disc
					   を、pole $i$ から pole $j$ に
					   \hfill} */
					/* {\ pole $k$ を経由して、移動する
					   \hfill} */
{
    if (n >= 1)
      {
	  move_discs(n - 1, i, k, j);	/* {\ 関数 {\tt move\_discs()}
					   の中で、さらに自分自身 \hfill} */
	  move_one_disc(i, j);		/* {\ {\tt move\_discs()} が使われ
					   ている。このような \hfill} */
	  print_result();		/* {\ 手法は、「再帰的呼びだし」
					   といわれる。一見 \hfill} */
	  move_discs(n - 1, k, j, i);	/* {\ 複雑に見える問題でも、
					   再帰的手法を用いると、\hfill} */
      }				/* {\ 簡単に解けてしまうことが、
					   しばしばある。\hfill} */
}

/* {\vskip 1cm
\special{epsfile=hanoi1.eps hscale=0.365 vscale=0.3}
\special{epsfile=hanoi2.eps hoffset=225 hscale=0.35 vscale=0.3}
\vskip 5cm

\noindent
たとえば、関数 {\tt move\_discs(4,i,j,k)} を呼び出すことは、
上図のような操作をすることに対応する。\hfill} */

int main(void)
{
    ptr[0] = disc[0] + ARRAY;
    ptr[1] = disc[1];
    ptr[2] = disc[2];

    init_array();
    move_discs(ARRAY, 0, 1, 2);		/* {\ {\tt ARRAY} 枚の disc を
					   pole 0 から pole 1 に pole 2
					   \hfill} */
					/* {\ を経由して、移動
					   する \hfill} */
    return 0;
}
