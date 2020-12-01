#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <readline/readline.h>
#include <sys/wait.h>
#include <unistd.h>

char **parse_input(char *);

int main() {
    char *input;
    char **command;
    pid_t child_pid;
    int stat_loc;

    while (1) {
        input = readline("unixsh> ");
        command = parse_input(input);

        if (!command[0]) {      /* Handle empty commands */
            free(input);
            free(command);
            continue;
        }

        child_pid = fork();     /* Create a child process */

        if (child_pid == 0) {       /* Execute the command in the child process */
            execvp(command[0], command);

        } else {       /* While the parent waits for the command to complete */
            waitpid(child_pid, &stat_loc, WUNTRACED);
        }

        free(input);
        free(command);
    }

    return 0;
}

char **parse_input(char *input) {
    /* Static memory allocation:
     Currently our command buffer allocates 8 blocks only, meaning it can properly parse an eight words size input */
    char **command = malloc(8 * sizeof(char *));

    char *separator = " ";
    char *parsed;
    int index = 0;

    parsed = strtok(input, separator);
    while (parsed != NULL) {
        command[index] = parsed;
        index++;

        parsed = strtok(NULL, separator);
    }

    command[index] = NULL;
    return command;
}
