
int fib(n) {
    int a = 1;
    int b = 1;
    while(n > 1) {
	int temp = a + b;
	a = b;
	b = temp;
	n--;
    }
    return a;
}

int main() {
    return fib(10);
}
