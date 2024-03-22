MOD = cpu_so1_1s2024
PWD :=$(shell pwd)
obj-m := $(MOD).o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean

run:
	sudo insmod $(MOD).ko

stop:
	sudo rmmod $(MOD)