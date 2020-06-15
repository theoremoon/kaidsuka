#include <cstdio>
#include <map>
#include <set>
#include <string>
#include <vector>

#include <retdec/retdec/retdec.h>

int main(int argc, char *argv[]) {
  if (argc == 2) {
    std::string input = argv[1];
    retdec::common::FunctionSet fs;
    auto res = retdec::disassemble(input, &fs);

    for (auto &f : fs) {
      std::printf("%s - 0x%x\n", f.getName().c_str(), f.getStart());
    }
  } else if (argc == 3) {
    std::string input = argv[1];
    std::string funcName = argv[2];
    retdec::common::FunctionSet fs;
    auto res = retdec::disassemble(input, &fs);

    const retdec::common::Function *func;
    for (auto &f : fs) {
      if (f.getName() == funcName) {
        func = &f;
        break;
      }
    }
    if (func == nullptr) {
      std::exit(1);
    }

    // pick destinations
    csh handle;
    cs_open(CS_ARCH_X86, CS_MODE_64, &handle);
    cs_option(handle, CS_OPT_DETAIL, CS_OPT_ON);

    for (auto &bb : func->basicBlocks) {
      std::printf("------\n");
      cs_regs regs_read, regs_write;
      uint8_t read_count, write_count;
      std::set<uint32_t> dests;
      for (auto *insn : bb.instructions) {
        cs_regs_access(handle, insn, regs_read, &read_count, regs_write,
                       &write_count);
        for (int i = 0; i < write_count; i++) {
          dests.insert(regs_write[i]);
        }
      }

      std::map<uint32_t, std::vector<std::string>> strands;
      for (auto &dest : dests) {
        strands[dest].resize(0);
      }

      for (auto *insn : bb.instructions) {
        cs_regs_access(handle, insn, regs_read, &read_count, regs_write,
                       &write_count);

        std::set<uint32_t> used;
        for (int i = 0; i < read_count; i++) {
          if (dests.find(regs_read[i]) != dests.end() &&
              used.find(regs_read[i]) == used.end()) {
            used.insert(regs_read[i]);
            strands[regs_read[i]].push_back(std::string(insn->mnemonic) + " " +
                                            std::string(insn->op_str));
          }
        }
        for (int i = 0; i < write_count; i++) {
          if (dests.find(regs_write[i]) != dests.end() &&
              used.find(regs_write[i]) == used.end()) {
            used.insert(regs_write[i]);
            strands[regs_write[i]].push_back(std::string(insn->mnemonic) + " " +
                                             std::string(insn->op_str));
          }
        }
      }

      for (auto &[reg, strand] : strands) {
        printf("REG: %s(%d)\n", cs_reg_name(handle, reg), reg);
        for (auto &line : strand) {
          printf("  %s\n", line.c_str());
        }
        printf("\n");
      }
    }
  } else {
    std::printf("Usage: %s <ELF> - Show function names\n", argv[0]);
    std::printf("Usage: %s <ELF> <FuncName> - Lift up function to llvm\n",
                argv[0]);
    return 1;
  }
  return 0;
}
