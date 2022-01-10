#include <iostream>

int main()
{
    int a = 42;

    for (int i = 0; i < 42; i++) {
        a += i;
    }

    std::cout << "Hello World " << a << std::endl;
    return 0;
}
