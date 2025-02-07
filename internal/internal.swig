%module internal 
%include "stdint.i"
%include <typemaps.i>

%ignore "lnf_mem_read_raw";
%ignore "lnf_mem_read_reset";
%ignore "lnf_error";

%{
#include <string.h>
#include <stdlib.h>
%}

%rename("open") "lnf_open";
%rename("read") "lnf_read";
%rename("write") "lnf_write";
%rename("info") "lnf_info";
%rename("close") "lnf_close";

%rename("rec_init") "lnf_rec_init";
%rename("rec_clear") "lnf_rec_clear";
%rename("rec_copy") "lnf_rec_copy";
%rename("rec_fset") "lnf_rec_fset";
%rename("rec_fget") "lnf_rec_fget";
%rename("rec_get_raw") "lnf_rec_get_raw";
%rename("rec_set_raw") "lnf_rec_set_raw";
%rename("rec_free") "lnf_rec_free";

%rename("filter_init") "lnf_filter_init";
%rename("filter_init_v1") "lnf_filter_init_v1";
%rename("filter_init_v2") "lnf_filter_init_v2";
%rename("filter_ffilter_ptr") "lnf_filter_ffilter_ptr";
%rename("filter_match") "lnf_filter_match";
%rename("filter_free") "lnf_filter_free";

%rename("mem_init") "lnf_mem_init";
%rename("mem_setopt") "lnf_mem_setopt";
%rename("mem_fadd") "lnf_mem_fadd";
%rename("mem_fastaggr") "lnf_mem_fastaggr";
%rename("mem_write") "lnf_mem_write";
%rename("mem_write_raw") "lnf_mem_write_raw";
%rename("mem_merge_threads") "lnf_mem_merge_threads";
%rename("mem_read") "lnf_mem_read";
%rename("mem_first_c") "lnf_mem_first_c";
%rename("mem_next_c") "lnf_mem_next_c";
%rename("mem_lookup_c") "lnf_mem_lookup_c";
%rename("mem_lookup_raw_c") "lnf_mem_lookup_raw_c";
%rename("mem_read_c") "lnf_mem_read_c";
%rename("mem_read_raw_c") "lnf_mem_read_raw_c";
%rename("mem_clean") "lnf_mem_clean";
%rename("mem_free") "lnf_mem_free";

%rename("ring_init") "lnf_ring_init";
%rename("ring_info") "lnf_ring_info";
%rename("ring_read") "lnf_ring_read";
%rename("ring_write") "lnf_ring_write";
%rename("ring_free") "lnf_ring_free";

%rename("fld_type") "lnf_fld_type";
%rename("fld_info") "lnf_fld_info";
%rename("fld_parse") "lnf_fld_parse";

%typemap(gotype) uint32_t "uint32"
%typemap(gotype) uint8_t "uint8"
%typemap(gotype) uint32_t[ANY] "*uint32"
%typemap(gotype) uint8_t[ANY] "*uint8"


%typemap(out) lnf_ip_t {
    $result = LnfIP{
        Data: [4]uint32{
            uint32($1.data[0]),
            uint32($1.data[1]),
            uint32($1.data[2]),
            uint32($1.data[3]),
        },
    };
}

%typemap(in) lnf_ip_t {
    static lnf_ip_t temp;
    temp.data[0] = (uint32_t)$input.Data[0];
    temp.data[1] = (uint32_t)$input.Data[1];
    temp.data[2] = (uint32_t)$input.Data[2];
    temp.data[3] = (uint32_t)$input.Data[3];
    $1 = &temp;
}
%rename(Lnf_ip_t) ""; 

%rename("%(strip:[LNF_])s") "";
%{
    #include "libnf.h"
%}

%include "libnf.h"
