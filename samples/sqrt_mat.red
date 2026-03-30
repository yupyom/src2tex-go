% {\bf sqrt\_mat.red}

% {\ For a given $2\times2$ matrix $\,A\ge0\,$, we shall prove that
% $$
% \sqrt{A}={1\over\sqrt{2\sqrt{ac-b^2}+a+c}}
% \pmatrix{
% \sqrt{ac-b^2}+a & b \cr
% b & \sqrt{ac-b^2}+c \cr}
% $$}
%

on div;

% {\ For a given matrix
% $
% A=\pmatrix{ a & b \cr b & c \cr}
% $
% ,\ we shall find a matrix
% $
% X=\pmatrix{ x & y \cr y & z \cr}
% $
% satisfying $A=X^2$.}
%

mat_a:=mat((a,b),(b,c));
mat_x:=mat((x,y),(y,z));

% {\ The equation
% $\,A=X^2\,$ is translated into algebraic equations
% $$
% x^2+y^2=a,\ y(x+z)=b,\ y^2+z^2=c \leqno(1)
% $$}
%

mat_xx:=mat_x*mat_x;

% {\ (1) is equivalent to
% $$
% x=\sqrt{a-y^2},
%	\ y\bigl(\sqrt{a-y^2}+\sqrt{c-y^2}\,\bigr)=b,
%	\ z=\sqrt{c-y^2} \leqno(2)
% $$
% We put
% $$
% f(y)=\sqrt{a-y^2},
% \ g(y)=y\bigl(\sqrt{a-y^2}+\sqrt{c-y^2}\,\bigr)-b,
% \ h(y)=\sqrt{c-y^2} \leqno(3)
% $$}
%

func_f:=rhs(first(solve(mat_xx(1,1)=a,x)));
func_h:=rhs(first(solve(mat_xx(2,2)=c,z)));
func_g:=sub({x=func_f,z=func_h},mat_xx(1,2))-b;

% {\ We define a function $\,g_1(w)\,$ by
% $$
% g_1(w)=g(y)\;{\sqrt{a-y^2}-\sqrt{c-y^2}\over y}\;\Bigg|_{y=1/w}
%	=-b\sqrt{aw^2-1}+b\sqrt{cw^2-1}+a-c \leqno(4)
% $$}
%

func_g1:=sub(y=1/w,func_g*(sqrt(a-y^2)-sqrt(c-y^2))/y);

% {\ We put
% $$
% g_2(v)
%	=g_1(w)\Big|_{w=\sqrt{(v^2+1)/a}}
%	=b\sqrt{-a+cv^2+c\over a}-bv+a-c \leqno(5)
% $$}
%

func_g2:=sub(w=sqrt((v^2+1)/a),func_g1);

% {\ Since
% $$
% \eqalign{
% g_2(v)&=0\cr
% g_2(v)+bv-a+c&=bv-a+c\cr
% (g_2(v)+bv-a+c)^2&=(bv-a+c)^2\cr
% {b^2\over a\ }\;(-a+cv^2+c)&=(bv-a+c)^2\cr
% }
% $$
% we put
% $$
% g_3(v)={b^2\over a\ }\;(-a+cv^2+c)-(bv-a+c)^2 \leqno(6)
% $$
% We can solve the equation $\,g_3(v)=0\,$ and get two solutions
% $$
% v_1={a+\sqrt{ac-b^2}\over b}\ ,\quad v_2={a-\sqrt{ac-b^2}\over b}
% $$
% Here we note that since
% $$
% g_2(v_2)=\sqrt{-2c\sqrt{ac-b^2}+ac-b^2+c^2}+\sqrt{ac-b^2}-c
% $$
% $\,g_2(v_2)\,$ is not always equal to 0. So, we can reject $\,v_2\,$.}
%

func_g3:=(func_g2+b*v-a+c)^2-(b*v-a+c)^2;
solutions:=solve(func_g3,v);
v1:=rhs(first(solutions));
v2:=rhs(second(solutions));
sub(v=v1,func_g2);
sub(v=v2,func_g2);

% {\ Since
% $$
% y=\sqrt{a\over v^2+1}
% $$
% we have
% $$
% y_1={b\over \sqrt{2\sqrt{ac-b^2}+a+c}}
% $$}
%

y1:=sqrt(a/(v1^2+1));

% {\ Hence, we obtain
% $$
% x_1=f(y_1)\ ,\quad z_1=h(y_1)
% $$
% this gives
% $$
% X_1
% =\pmatrix{x_1 & y_1 \cr
%	y_1 & z_1 \cr}
% ={1\over d_1}\pmatrix{ \sqrt{ac-b^2}+a & b \cr
%	b & \sqrt{ac-b^2}+c \cr}
% $$
% where $\,d_1=\sqrt{2\sqrt{ac-b^2}+a+c}\,$.
% Direct computation shows $\,X_1^2=A\,$.}
%

x1:=sub(y=y1,func_f);
z1:=sub(y=y1,func_h);
mat_x1:=mat((x1,y1),(y1,z1));
d1:=sqrt(2*sqrt(a*c - b**2) + a + c);
d1*mat_x1;
mat_x1*mat_x1;

end;
