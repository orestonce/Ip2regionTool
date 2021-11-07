#include "txttodbform.h"
#include "ui_txttodbform.h"
#include "ip2region.h"
#include <QFileDialog>
#include <QMessageBox>

TxtToDbForm::TxtToDbForm(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::TxtToDbForm)
{
    ui->setupUi(this);
    this->refresh_startConvert_Enable();
}

TxtToDbForm::~TxtToDbForm()
{
    delete ui;
}

void TxtToDbForm::on_pushButton_selectDb_clicked()
{
    QString dbFileNameStr = QFileDialog::getSaveFileName(this, ui->label_selectDb->text(),  "", "*.db");
    ui->lineEdit_selectDb->setText(dbFileNameStr);
    this->refresh_startConvert_Enable();
}


void TxtToDbForm::on_pushButton_selectTxt_clicked()
{
    QString dbFileNameStr = QFileDialog::getOpenFileName(this, ui->label_selectDb->text(),  "", "*.txt");
    ui->lineEdit_selectTxt->setText(dbFileNameStr);
    this->refresh_startConvert_Enable();
}

void TxtToDbForm::on_pushButton_startConvert_clicked()
{
    std::string dbFileName = ui->lineEdit_selectDb->text().toStdString();
    std::string txtFileName = ui->lineEdit_selectTxt->text().toStdString();
    std::string errMsg = ConvertTxtToDb(txtFileName, dbFileName);
    if (!errMsg.empty()) {
        QMessageBox::about(this, "错误", errMsg.c_str());
        return;
    }
    QMessageBox::about(this, "成功", "转换成功!");
}

void TxtToDbForm::refresh_startConvert_Enable()
{
    if (ui->lineEdit_selectDb->text().isEmpty() || ui->lineEdit_selectTxt->text().isEmpty()) {
        ui->pushButton_startConvert->setEnabled(false);
    } else {
        ui->pushButton_startConvert->setEnabled(true);
    }
}
