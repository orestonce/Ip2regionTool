#pragma once

#include <string>
#include <vector>
#include <cstdint>
#include <map>
//Qt Creator 需要在xxx.pro 内部增加静态库的链接声明
//LIBS += -L$$PWD -lip2region-impl

struct ConvertDbToTxt_Req{
	std::string DbFileName;
	std::string TxtFileName;
	bool Merge;
	ConvertDbToTxt_Req(): Merge(false){}
};
std::string ConvertDbToTxt(ConvertDbToTxt_Req in0);
struct ConvertTxtToDb_Req{
	std::string TxtFileName;
	std::string DbFileName;
	std::string RegionCsvFileName;
	bool Merge;
	ConvertTxtToDb_Req(): Merge(false){}
};
std::string ConvertTxtToDb(ConvertTxtToDb_Req in0);
