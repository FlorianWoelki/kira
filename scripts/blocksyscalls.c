/*
 * This program is a simple script that disables the socket syscall.
 */
#include <stdio.h>
#include <errno.h>
#include <unistd.h>
#include <sys/prctl.h>
#include <seccomp.h>

int main(int argc, char *argv[])
{
    // Disable any new capabilities.
    prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0);

    // Initialize the filter.
    scmp_filter_ctx ctx = seccomp_init(SCMP_ACT_ALLOW);
    if (ctx == NULL) {
        perror("seccomp_init failed");
        return 1;
    }

    // Add the rule to deny the socket syscall.
    if (seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EPERM), SCMP_SYS(socket), 0) == -1) {
        perror("seccomp_rule_add failed");
        return 1;
    }

    // Load the filter.
    if (seccomp_load(ctx) == -1) {
        perror("seccomp_load failed");
        return 1;
    }

    // Execute the arguments that were passed to this program.
    execvp(argv[1], &argv[1]);
    return 0;
}
