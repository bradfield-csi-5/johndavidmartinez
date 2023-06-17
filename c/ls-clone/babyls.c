#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <dirent.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pwd.h>
#include <grp.h>
#include <time.h>

#define MAX_DIRECTORY_SIZE 100

// next implement long format
// stretch implement colors

// permissions mask values for st_mode organized
// in order of long format display
int permission_lookup[] = {
	S_IRUSR, S_IWUSR, S_IXUSR,
	S_IRGRP, S_IWGRP, S_IXGRP,
	S_IROTH, S_IWOTH, S_IXOTH
};

// character lookup table for permission_lookup masks
char permission_char[] = {
	'r', 'w', 'x',
	'r', 'w', 'x',
	'r', 'w', 'x'
};

// file type masks
// saw this after making this table
// but man 7 inode shows you can do pre-made macros like S_ISDIR(m)
int file_type_lookup[] = {
	S_IFMT, S_IFSOCK, S_IFLNK,
	S_IFREG, S_IFBLK, S_IFDIR,
	S_IFCHR, S_IFIFO
};

// file type characters
// should add support beyond directory
char file_type_char[] = {
	'-', '-', '-',
	'-', '-', 'd',
	'-', '-'
};

int alpha_dirbuf_comp(const void *d1, const void *d2) {
	struct dirent *dirbuf1 = (struct dirent *)d1;
	struct dirent *dirbuf2 = (struct dirent *)d2;
	return strcmp(dirbuf1->d_name, dirbuf2->d_name);
}

// could be generalized to all comp funcs if new sorting is added
int r_alpha_dirbuf_comp(const void *d1, const void *d2) {
	int res = alpha_dirbuf_comp(d1, d2);
	return res * -1;
}

void print_directories_names(struct dirent dirarr[], int len) {
    	int i;
    	for (i = 0; i < len; i++) {
    		printf("%s  ", dirarr[i].d_name); 
    	}
    	printf("\n");
}

void permissions_str(unsigned long st_mode, char *buf) {
	int i;
	for (i = 0; i < 10; i++) {
		buf[i] = '-';
	}
	for (i = 0; i < 9; i++) {
		if ((st_mode & permission_lookup[i]) == permission_lookup[i]) {
			buf[i + 1] = permission_char[i];
		}
	}
	for (i = 0; i < 8; i++) {
		if ((st_mode & file_type_lookup[i]) == file_type_lookup[i]) {
			buf[0] = file_type_char[i];
		}
	}
	buf[10] = '\0';
}

void trimed_ts(long epoch, char *buf) {
	int i;
	char* full_ts_string = ctime(&epoch);
	for (i = 0; i < 12; i++) {
		buf[i] = full_ts_string[i + 4];
	}
	buf[12] = '\0';
}

void print_directories(struct dirent dirarr[], struct stat statarr[], 
		       struct passwd passwdarr[], struct group grouparr[], int len) {
	int i;
	unsigned long totalblks = 0;
        char permissions[11], trimed[13];

	for (i = 0; i < len; i++) {
		totalblks += statarr[i].st_blocks;
	}
	// divide by two is a HACK
	// I don't understand block size
	// blksize_t st_blksize;     /* Block size for filesystem I/O
        // blkcnt_t  st_blocks;      /* Number of 512B blocks allocated
	// Maybe you do math based on st_blksize?
	printf("total %ld\n", totalblks / 2);
	for (i = 0; i < len; i++) {
                permissions_str(statarr[i].st_mode, permissions);
                trimed_ts(statarr[i].st_mtim.tv_sec, trimed);
		// %5ld should be dynamic.. 
		printf("%s %ld %s %s %5ld %s %s\n",
				permissions,
				statarr[i].st_nlink,
				passwdarr[i].pw_name, grouparr[i].gr_name,
				statarr[i].st_size, 
				trimed,
				dirarr[i].d_name);
	}
}

// how people organize string manipulation in c is beyond me
// this would be difficult to manage for long running programs
char* dir_cat(char* base, char* path) {
	char* full_path = malloc(MAX_DIRECTORY_SIZE * sizeof(char));
	strncpy(full_path, base, MAX_DIRECTORY_SIZE);
	strncat(full_path, path, MAX_DIRECTORY_SIZE);
	return full_path;
}

int main(int argc, char *argv[]) {
	int ignorehidden = 1;
	int asort = 1;
	int color = 1;
	int reverse = 0;
	int long_list_fmt = 0;
	int i;
	// Process Arguments
        char* dirz_to_ls[4];
        int dirz_idx = 0;
        for (i = 0; i < 4; i++) {
            dirz_to_ls[i] = malloc(MAX_DIRECTORY_SIZE * sizeof(char));
        }
        // default CWD
        dirz_to_ls[0][0] = '.';
	dirz_to_ls[0][1] = '\0';
	for (i = 1; i < argc; i++) {
		if (argv[i][0] != '-') {
			// assume non-flag argument is directory
                        strncpy(dirz_to_ls[dirz_idx], argv[i], MAX_DIRECTORY_SIZE);
                        dirz_idx++;
		}
		if (strcmp(argv[i], "-a") == 0) {
			ignorehidden = 0;
		}
		if (strcmp(argv[i], "-A") == 0) {
			ignorehidden = 2;
		}
		if (strcmp(argv[i], "-f") == 0) {
			asort = 0;
			ignorehidden = 0;
			color = 0;
		}
		if (strcmp(argv[i], "-r") == 0) {
			reverse = 1;
		}
		if (strcmp(argv[i], "-l") == 0) {
			long_list_fmt = 1;
		}
	}
	// add relative / if not provided
        for (i = 0; i < 4; i ++) {
        	char* dir_to_ls_ptr = dirz_to_ls[i];
        	while (*(dir_to_ls_ptr + 1) != '\0') {
        		dir_to_ls_ptr++;
        	}
        	if (*dir_to_ls_ptr != '/') {
        		dir_to_ls_ptr++;
        		*dir_to_ls_ptr = '/';
        		dir_to_ls_ptr++;
        		*dir_to_ls_ptr = '\0';
        	}
        }

	// Open main directory fd
	// Need to update this to take first argument as a path
	int fd = open(dirz_to_ls[0], O_RDONLY|O_NONBLOCK|O_CLOEXEC|O_DIRECTORY);
	if (fd == -1) {
		printf("oops");
		return 1;
	}
	// Read directory file and initialize buffers
	struct stat statbuf;
	fstat(fd, &statbuf);
	DIR *dirptr = fdopendir(fd);
	struct dirent *dirbuf;
	struct dirent dirarr[MAX_DIRECTORY_SIZE];
	int dirarridx = 0;
	// Read directory entries and store in array
	while (1) {
		dirbuf = readdir(dirptr);
		if (dirbuf == NULL || dirarridx == MAX_DIRECTORY_SIZE) {
			break;
		}
		// ignore hidden if required
		if (ignorehidden == 2) {
			if (dirbuf->d_name[0] == '.' && 
					(strcmp(dirbuf->d_name, ".") == 0 || strcmp(dirbuf->d_name, "..") == 0)) {
				continue;
			}
		}
		// ignore hidden if required
		if (ignorehidden == 1 && dirbuf->d_name[0] == '.') {
			continue;
		}
		dirarr[dirarridx] = *dirbuf;
		dirarridx++;
	}	
	// sort if required
	if (asort) {
		if (reverse) {
			//be fun to implement my own sort and might help me understand the bug I had
			//trying to call qsort with an array of struct dirent* instead of struct dirent
			qsort(dirarr, dirarridx, sizeof(struct dirent), r_alpha_dirbuf_comp);
		} else {
			qsort(dirarr, dirarridx, sizeof(struct dirent), alpha_dirbuf_comp);
		}
	}
        // close directory
        closedir(dirptr);
	// load additional data if required (long format flag)
	struct stat statarr[MAX_DIRECTORY_SIZE];
	struct passwd passwdarr[MAX_DIRECTORY_SIZE];
	struct group grouparr[MAX_DIRECTORY_SIZE];
	i = 0;
	while (long_list_fmt) {
		if (i == dirarridx) {
			break;
		}
		dirbuf = &dirarr[i];
		stat(dir_cat(dirz_to_ls[0], dirbuf->d_name), &statbuf);
		statarr[i] = statbuf;
		passwdarr[i] = *getpwuid(statbuf.st_uid);
		grouparr[i] = *getgrgid(statbuf.st_gid);

		i++;
	}
	if (long_list_fmt) {
	        print_directories(dirarr, statarr, passwdarr, grouparr, dirarridx);
	} else {
	        print_directories_names(dirarr, dirarridx);
	}
	close(fd);
        for (i = 0; i < 4; i++) {
            free(dirz_to_ls[i]);
        }
}
