#include <cstdio>
#include <cstdlib>
#include <string>
#include <vector>

#include <retdec/capstone2llvmir/capstone2llvmir.h>
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

    // find function
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

    // extract machine code bytes
    std::vector<uint8_t> code;
    for (auto &bb : func->basicBlocks) {
      for (auto *insn : bb.instructions) {
        for (int i = 0; i < insn->size; i++) {
          code.push_back(insn->bytes[i]);
        }
      }
    }

    llvm::LLVMContext ctx;
    llvm::Module module("test", ctx);
    auto *llvm_f = llvm::Function::Create(
        llvm::FunctionType::get(llvm::Type::getVoidTy(ctx), false),
        llvm::GlobalValue::ExternalLinkage, funcName.c_str(), &module);
    llvm::BasicBlock::Create(module.getContext(), "entry", llvm_f);
    llvm::IRBuilder<> irb(&llvm_f->front());
    try {
      auto c2l = retdec::capstone2llvmir::Capstone2LlvmIrTranslator::createArch(
          CS_ARCH_X86, &module, CS_MODE_64, CS_MODE_LITTLE_ENDIAN);
      c2l->translate(code.data(), code.size(), func->getStart(), irb);
      module.print(llvm::outs(), nullptr);
    } catch (...) {
      std::printf("ERR\n");
      std::exit(1);
    }

  } else {
    std::printf("Usage: %s <ELF> - Show function names\n", argv[0]);
    std::printf("Usage: %s <ELF> <FuncName> - Lift up function to llvm\n",
                argv[0]);
    return 1;
  }
  return 0;
}
