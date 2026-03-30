{ {\hrulefill }

  {\ % beginning of TeX mode } 

  {\ \centerline{\bf Towers of Hanoi (Pascal)} } 

  {\ \begin{quote} } 
  {\ \noindent } 
  {\ This program gives an answer to the following famous problem (towers of } 
  {\ Hanoi). } 
  {\ There is a legend that when one of the temples in Hanoi was constructed, } 
  {\ three poles were erected and a tower consisting of 64 golden discs was } 
  {\ arranged on one pole, their sizes decreasing regularly from bottom to top. } 
  {\ The monks were to move the tower of discs to the opposite pole, moving } 
  {\ only one at a time, and never putting any size disc above a smaller one. } 
  {\ The job was to be done in the minimum numbers of moves. What strategy for } 
  {\ moving discs will accomplish this optimum transfer? } 
  {\ \end{quote} } 

  {\ % end of TeX mode } 

  {\hrulefill } }

program Hanoi;

const
  ARRAYSIZE = 8;                                  { disc の数 }

var
  disc: array[0..2, 0..ARRAYSIZE-1] of integer;   { disc に関するデータの置き場所 }
  ptr: array[0..2] of integer;                    { スタックのポインタ }
  counter: integer;                               { 移動回数カウンタ }

procedure InitArray;                              { disc 移動用アレイ }
var
  j: integer;
begin
  for j := 0 to ARRAYSIZE - 1 do
  begin
    disc[0][j] := ARRAYSIZE - j;
    disc[1][j] := 0;
    disc[2][j] := 0;
  end;
  ptr[0] := ARRAYSIZE;
  ptr[1] := 0;
  ptr[2] := 0;
end;

procedure PrintResult;                            { 結果の表示 }
var
  i, j: integer;
begin
  counter := counter + 1;
  writeln('---', counter, '---');
  for i := 0 to 2 do
  begin
    write('[', i, '] ');
    for j := 0 to ARRAYSIZE - 1 do
    begin
      if disc[i][j] <> 0 then
        write(disc[i][j], ' ')
      else
        break;
    end;
    writeln;
  end;
end;

procedure MoveOneDisc(i, j: integer);       { 1枚の disc を pole i から pole j に移動する }
begin
  ptr[i] := ptr[i] - 1;
  disc[j][ptr[j]] := disc[i][ptr[i]];
  ptr[j] := ptr[j] + 1;
  disc[i][ptr[i]] := 0;
end;

procedure MoveDiscs(n, i, j, k: integer);   { n 枚を移動 }
begin                                       { pole i から pole j }
  if n >= 1 then                            { pole k 経由で }
  begin
    MoveDiscs(n - 1, i, k, j);
    MoveOneDisc(i, j);                      { move_discs() で自分自身を移動 }
    PrintResult;
    MoveDiscs(n - 1, k, j, i);
  end;
end;

begin
  counter := 0;
  InitArray;
  PrintResult;
  MoveDiscs(ARRAYSIZE, 0, 2, 1);
end.
