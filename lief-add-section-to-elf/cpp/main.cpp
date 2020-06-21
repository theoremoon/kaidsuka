#include <LIEF/LIEF.hpp>
#include <iostream>
#include <string>

using namespace LIEF;

int main(int argc, char **argv) {
  try {
    std::vector<uint8_t> data(100, 0); // really?

    auto binary = ELF::Parser::parse(argv[1]);
    auto header = binary->header();

    // TODO confirm it is stripped and static

    // add sections .symtab and .strtab
    ELF::Section symtab{".symtab"};
    symtab.type(ELF::ELF_SECTION_TYPES::SHT_SYMTAB);
    symtab.alignment(8);
    symtab.entry_size(0x18);
    symtab.link(header.numberof_sections() + 1);
    symtab.content(data);

    ELF::Section strtab{".strtab"};
    strtab.type(ELF::ELF_SECTION_TYPES::SHT_STRTAB);
    strtab.alignment(8);
    strtab.entry_size(1);
    strtab.content(data);

    binary->add(symtab, false); // dont load at exection
    binary->add(strtab, false);

    // idk what is it
    {
      ELF::Symbol symbol{""};
      symbol.type(ELF::ELF_SYMBOL_TYPES::STT_NOTYPE);
      symbol.value(0);
      symbol.binding(ELF::SYMBOL_BINDINGS::STB_LOCAL);
      symbol.size(0);
      symbol.shndx(0);

      binary->add_static_symbol(symbol);
    }

    // create symbols
    ELF::Symbol symbol{argv[2]};
    symbol.type(ELF::ELF_SYMBOL_TYPES::STT_FUNC);
    symbol.value(std::stoul(argv[3], nullptr, 16));
    symbol.binding(ELF::SYMBOL_BINDINGS::STB_LOCAL);
    symbol.shndx(1);

    // add symbol to binary
    binary->add_static_symbol(symbol);

    // write out to file
    binary->write(argv[4]);
  } catch (const LIEF::exception &err) {
    std::cerr << err.what() << std::endl;
  }
  return 0;
}
