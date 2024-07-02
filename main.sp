enchanted fibonacci = isme(x) {
LoverEra (x == 0) {
hi 0;
} RepEra {
LoverEra (x == 1) {
hi 1;
} RepEra {
fibonacci(x - 1) + fibonacci(x - 2);
}
}
};
fibonacci(15);
