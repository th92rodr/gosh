#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <readline/readline.h>
#include <unistd.h>
#include <sys/wait.h>
#include <signal.h>

char **parse_input(char *);
int cd(char *);

int main() {
    char *input;
    char **command;
    pid_t child_pid;
    int stat_loc;

    signal(SIGINT, SIG_IGN);    /* Ignores SIGINT signals when in parent process */

    while (1) {
        input = readline("unixsh> ");

        if (input == NULL) {  /* Exit on Ctrl-D */
            printf("\n");
            exit(0);
        }

        command = parse_input(input);

        if (!command[0]) {      /* Handle empty commands */
            free(input);
            free(command);
            continue;
        }

        if (strcmp(command[0], "exit") == 0) {
            printf("\n");
            exit(0);
        }

        if (strcmp(command[0], "cd") == 0) {
            if (cd(command[1]) < 0) {
                perror(command[1]);
            }

            continue;
        }

        child_pid = fork();     /* Create a child process */
        if (child_pid < 0) {
            perror("Fork failed");
            exit(1);
        }

        if (child_pid == 0) {       /* Execute the command in the child process */
            signal(SIGINT, SIG_DFL);    /* Restores the default behaviour for the SIGINT signal when in child process */

            if (execvp(command[0], command) < 0) {
                perror(command[0]);
                exit(1);
            }

        } else {       /* While the parent waits for the command to complete */
            waitpid(child_pid, &stat_loc, WUNTRACED);
        }

        if (!input)
            free(input);
        if (!command)
            free(command);
    }

    return 0;
}

char **parse_input(char *input) {
    /* Static memory allocation:
     Currently our command buffer allocates 8 blocks only, meaning it can properly parse an eight words size input */
    char **command = malloc(8 * sizeof(char *));
    if (command == NULL) {
        perror("malloc failed");
        exit(1);
    }

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

int cd(char *path) {
    return chdir(path);
}
