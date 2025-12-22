reset

if (!exists("file")) file = "results/results_J_1_h_0.csv"

set datafile separator ";"
set term pngcairo enhanced size 800,600

# Энергия U(T)
set xlabel "T"
set ylabel "<E>/N"
set grid
set key top left
# В Windows путь в переменной file может содержать обратные слэши,
# поэтому имя вида "E_\<path>.png" создаёт несуществующую папку "E_<...>".
# Вместо этого сохраняем файлы рядом с CSV: "<file>_E.png", "<file>_C.png", "<file>_M.png".
set output sprintf("%s_E.png", file)
plot file using 1:2 with lines lw 2 title "U(T)"

# Теплоёмкость C(T)
set ylabel "C/N"
set output sprintf("%s_C.png", file)
plot file using 1:4 with lines lw 2 title "C(T)"

# Магнитизация m(T)
set ylabel "m"
set output sprintf("%s_M.png", file)
plot file using 1:3 with lines lw 2 title "m(T)"


