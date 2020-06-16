#include <cstdio>
#include <map>
#include <set>
#include <string>
#include <utility>
#include <vector>

#include <llvm/IR/InstIterator.h>
#include <retdec/capstone2llvmir/capstone2llvmir.h>
#include <retdec/retdec/retdec.h>

void convertInsts2LLVMIRs(llvm::Module &llvm_module,
                          const std::vector<uint8_t> &code,
                          const char *f_name) {
  auto *llvm_f = llvm::Function::Create(
      llvm::FunctionType::get(llvm::Type::getVoidTy(llvm_module.getContext()),
                              false),
      llvm::GlobalValue::ExternalLinkage, f_name, &llvm_module);
  llvm::BasicBlock::Create(llvm_module.getContext(), f_name, llvm_f);
  llvm::IRBuilder<> irb(&llvm_f->front());

  // throwable
  auto c2l = retdec::capstone2llvmir::Capstone2LlvmIrTranslator::createArch(
      CS_ARCH_X86, &llvm_module, CS_MODE_64, CS_MODE_LITTLE_ENDIAN);
  c2l->translate(code.data(), code.size(), 0x1000, irb);
}

std::map<uint32_t, std::vector<uint8_t>>
basicblockToStrands(csh handle, const retdec::common::BasicBlock &bb) {
  cs_regs regs_read, regs_write;
  uint8_t read_count, write_count;

  // list destination registers
  std::set<uint32_t> dests;
  for (const auto &insn : bb.instructions) {
    cs_regs_access(handle, insn, regs_read, &read_count, regs_write,
                   &write_count);
    for (int i = 0; i < write_count; i++) {
      dests.insert(regs_write[i]);
    }
  }

  // key: register, value: strand
  std::map<uint32_t, std::vector<uint8_t>> strands;
  for (auto &dest : dests) {
    strands[dest].resize(0);
  }

  for (const auto &insn : bb.instructions) {
    cs_regs_access(handle, insn, regs_read, &read_count, regs_write,
                   &write_count);

    std::set<uint32_t> used;
    for (int i = 0; i < read_count; i++) {
      if (dests.find(regs_read[i]) != dests.end() &&
          used.find(regs_read[i]) == used.end()) {
        used.insert(regs_read[i]);
        for (int j = 0; j < insn->size; j++) {
          strands[regs_read[i]].push_back(insn->bytes[j]);
        }
      }
    }
    for (int i = 0; i < write_count; i++) {
      if (dests.find(regs_write[i]) != dests.end() &&
          used.find(regs_write[i]) == used.end()) {
        for (int j = 0; j < insn->size; j++) {
          strands[regs_read[i]].push_back(insn->bytes[j]);
        }
      }
    }
  }
  return std::move(strands);
}

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
      std::printf("------------------\n");
      auto strands = basicblockToStrands(handle, bb);
      for (auto &[reg, strand] : strands) {
        try {
          llvm::LLVMContext ctx;
          llvm::Module llvm_module("test", ctx);
          convertInsts2LLVMIRs(llvm_module, strand, "FUNC_XXXXX");
          std::vector<llvm::Instruction *> strand_instructions;
          for (auto &func : llvm_module.getFunctionList()) {
            for (auto it = llvm::inst_begin(func); it != llvm::inst_end(func);
                 it++) {
              strand_instructions.push_back(&*it);
            }
          }
          std::printf("%d\n", strand_instructions.size());

          std::string ir_buf;
          llvm::raw_string_ostream rso(ir_buf);
          for (auto &ir : strand_instructions) {
            ir->print(rso);
            ir_buf.append("\n");
          }
          std::printf("%s\n", ir_buf.c_str());
        } catch (...) {
          std::printf("FAILED TO COVNVERT TO LLVM IR\n");
        }
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
