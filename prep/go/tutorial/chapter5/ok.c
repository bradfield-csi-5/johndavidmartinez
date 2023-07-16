#include <dirent.h>
#include <grp.h>
#include <inttypes.h>
#include <pwd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>
#include <unistd.h>

void fileout(char *filename, struct stat fs);

int main(int argc, char *argv[])
{
    char *base;
    DIR *dp;
    struct dirent *entry;
    struct stat fs;
    int r;
    char *filename;

    if (argc < 2) {
        base = ".";
    }
    else {
        base = argv[1];
    }

    r = stat(base, &fs);
    if (r == -1) {
        perror(argv[0]);
        return 1;
    }
    if (!S_ISDIR(fs.st_mode)) {
        fileout(base, fs);
        return 0;
    }

    dp = opendir(base);
    while ((entry = readdir(dp)) != NULL) {
        filename = entry->d_name;
        r = stat(filename, &fs);
        if (r == -1) {
            perror(argv[0]);
            closedir(dp);
            return 1;
        }
        fileout(filename, fs);
    }

    closedir(dp);
    return 0;
}

void mode_str(mode_t stm, char *str);
void modified_ts_str(time_t *t, char *str);

#define TS_SIZE 13
#define MODE_SIZE 11

void fileout(char *filename, struct stat fs)
{
    char mode[MODE_SIZE];
    mode_str(fs.st_mode, mode);

    char timestamp[TS_SIZE];
    modified_ts_str(&fs.st_atime, timestamp);

    struct passwd *pws;
    pws = getpwuid(fs.st_uid);

    struct group *g;
    g = getgrgid(fs.st_gid);

    printf("%s %s %s %8" PRIdMAX " %s %s\n", mode, pws->pw_name, g->gr_name,
           (intmax_t)fs.st_size, timestamp, filename);
}

void modified_ts_str(time_t *t, char *str)
{
    struct tm *clock = localtime(t);
    strftime(str, TS_SIZE, "%b %d %R", clock);
}

void mode_str(mode_t stm, char *str)
{
    if (S_ISDIR(stm)) {
        *str = 'd';
    }
    else if (S_ISLNK(stm)) {
        *str = 'l';
    }
    else if (S_ISBLK(stm)) {
        *str = 'b';
    }
    else if (S_ISCHR(stm)) {
        *str = 'c';
    }
    else if (S_ISFIFO(stm)) {
        *str = 'p';
    }
    else if (S_ISSOCK(stm)) {
        *str = 's';
    }

#ifdef S_ISWHT
    else if (S_ISWHT(stm)) {
        *str = 'w';
    }
#endif

    else {
        *str = '-';
    }

    *(++str) = stm & S_IRUSR ? 'r' : '-';
    *(++str) = stm & S_IWUSR ? 'w' : '-';
    *(++str) = stm & S_IXUSR ? 'x' : '-';
    *(++str) = stm & S_IRGRP ? 'r' : '-';
    *(++str) = stm & S_IWGRP ? 'w' : '-';
    *(++str) = stm & S_IXGRP ? 'x' : '-';
    *(++str) = stm & S_IROTH ? 'r' : '-';
    *(++str) = stm & S_IWOTH ? 'w' : '-';
    *(++str) = stm & S_IXOTH ? 'x' : '-';
    *(++str) = '\0';
}
