#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/sysinfo.h>
#include <linux/mm.h>
#include <linux/seq_file.h>

// Obteniendo estadísticas del sistema
struct sysinfo si;

static void init_meminfo(void) {
    si_meminfo(&si);
}

static int escribir_archivo(struct seq_file *archivo, void *v) {
    init_meminfo();
    // total en uso
    seq_printf(archivo, "%lu\n", si.totalram - si.freeram);
    // ram total
    seq_printf(archivo, "%lu\n", si.totalram);
    return 0;
}

static int al_abrir(struct inode *inode, struct file *file) {
    return single_open(file, escribir_archivo, NULL);
}

static struct proc_ops operaciones = {
    .proc_open = al_abrir, 
    .proc_read = seq_read
};

static int __init ram_init(void) {
    proc_create("ram_so1_1s2024", 0, NULL, &operaciones);
    printk(KERN_INFO "Iniciando módulo de RAM\n");
    return 0;
}

static void __exit ram_exit(void) {
    remove_proc_entry("ram_so1_1s2024", NULL);
    printk(KERN_INFO "Saliendo del módulo de RAM\n");
}

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Módulo de RAM Proyecto 1 - 201700761");
MODULE_AUTHOR("Elder Anibal Pum Rojas");

module_init(ram_init);
module_exit(ram_exit);