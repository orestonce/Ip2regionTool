#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QFileDialog>
#include <QMessageBox>
#include <QDebug>
#include "ip2region.h"

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    this->setTxtToXdb_IsRuning(false);
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_pushButton_input_db_clicked()
{
    QString input = QFileDialog::getOpenFileName(this, "",  "", "*.db *.xdb");
    if (!input.isEmpty()) {
        ui->lineEdit_input_db->setText(input);
        int version = GetDbVersionByName(input.toStdString());
        if (version==1) {
            ui->radioButton_dbv1->setChecked(true);
        } else {
            ui->radioButton_dbv2->setChecked(true);
        }
    }
}


void MainWindow::on_pushButton_output_txt_clicked()
{
    QString output = QFileDialog::getSaveFileName(this, "",  "", "*.txt");
    if (!output.isEmpty()) {
        ui->lineEdit_output_txt->setText(output);
    }
}

void MainWindow::on_pushButton_input_txt_clicked()
{
    QString input = QFileDialog::getOpenFileName(this, "",  "", "*.txt");
    if (!input.isEmpty()) {
        ui->lineEdit_input_txt->setText(input);
    }
}

void MainWindow::on_pushButton_output_db_clicked()
{
    QString output = QFileDialog::getSaveFileName(this, "",  "", "*.db");
    if (!output.isEmpty()) {
        ui->lineEdit_output_db->setText(output);
    }
}

void MainWindow::on_pushButton_input_regin_csv_clicked()
{
    QString input = QFileDialog::getOpenFileName(this, "", "", "*.csv");
    if (!input.isEmpty()) {
        ui->lineEdit_input_regin_csv->setText(input);
    }
}

void MainWindow::on_pushButton_DbToTxt_clicked()
{
    QString db = ui->lineEdit_input_db->text();
    QString txt = ui->lineEdit_output_txt->text();
    if (db.isEmpty()||txt.isEmpty()) {
        return;
    }
    ConvertDbToTxt_Req req;
    if (ui->radioButton_dbv1->isChecked()) {
        req.DbVersion = 1;
    } else {
        req.DbVersion = 2;
    }
    req.DbFileName = db.toStdString();
    req.TxtFileName = txt.toStdString();
    req.Merge = ui->checkBox_DbToTxt_merge->isChecked();
    std::string errMsg = ConvertDbToTxt(req);
    if (errMsg.empty()) {
        Toast::Instance()->SetSuccess("转换成功!");
        return;
    }
    Toast::Instance()->SetError(QString::fromStdString(errMsg));
}

void MainWindow::on_pushButton_TxtToDb_clicked()
{
    QString txt = ui->lineEdit_input_txt->text();
    QString db = ui->lineEdit_output_db->text();
    if (db.isEmpty()||txt.isEmpty()) {
        return;
    }
    ConvertTxtToDb_Req req;
    req.TxtFileName = txt.toStdString();
    req.DbFileName = db.toStdString();
    req.RegionCsvFileName = ui->lineEdit_input_regin_csv->text().toStdString();
    req.Merge = ui->checkBox_TxtToDb_merge->isChecked();
    std::string errMsg = ConvertTxtToDb(req);
    if (errMsg.empty()) {
        Toast::Instance()->SetSuccess("转换成功!");
        return;
    }
    Toast::Instance()->SetError(QString::fromStdString(errMsg));
}

void MainWindow::on_tabWidget_currentChanged(int index)
{
    ui->lineEdit_input_db->clear();
    ui->lineEdit_input_txt->clear();
    ui->lineEdit_output_db->clear();
    ui->lineEdit_output_txt->clear();
    ui->lineEdit_input_regin_csv->clear();
    ui->checkBox_DbToTxt_merge->setChecked(false);
    ui->checkBox_TxtToDb_merge->setChecked(false);
}


void MainWindow::on_pushButton_xdb_srcFile_clicked()
{
    QString input = QFileDialog::getOpenFileName(this, "", "", "*.txt");
    if (!input.isEmpty()) {
        ui->lineEdit_xdb_srcFile->setText(input);
    }
}

void MainWindow::on_pushButton_xdb_dstFile_clicked()
{
    QString output = QFileDialog::getSaveFileName(this, "",  "", "*.xdb");
    if (!output.isEmpty()) {
        ui->lineEdit_xdb_dstFile->setText(output);
    }
}

void MainWindow::on_pushButton_xdb_start_clicked()
{
    TxtToXdbReq req;
    req.SrcFile = ui->lineEdit_xdb_srcFile->text().toStdString();
    req.DstFile = ui->lineEdit_xdb_dstFile->text().toStdString();
    if (req.SrcFile.empty() || req.DstFile.empty() || ui->lineEdit_xdb_srcFile->isEnabled() == false) {
        return;
    }
    this->setTxtToXdb_IsRuning(true);

    req.IndexPolicyS = ui->comboBox_xdb_indexPolicy->currentText().toStdString();
    m_syncUi.AddRunFnOn_OtherThread([=](){
        std::string errMsg = TxtToXdb(req);
        m_syncUi.AddRunFnOn_UiThread([=](){
            this->setTxtToXdb_IsRuning(false);
            if (errMsg.empty()) {
                Toast::Instance()->SetSuccess("转换成功");
                return;
            }
            Toast::Instance()->SetError("失败 " + QString::fromStdString(errMsg));
        });
    });
}

void MainWindow::setTxtToXdb_IsRuning(bool runing)
{
    ui->lineEdit_xdb_srcFile->setEnabled(runing == false);
    ui->lineEdit_xdb_dstFile->setEnabled(runing == false);
    ui->comboBox_xdb_indexPolicy->setEnabled(runing == false);
    ui->pushButton_xdb_dstFile->setEnabled(runing == false);
    ui->pushButton_xdb_srcFile->setEnabled(runing == false);
    ui->pushButton_xdb_start->setEnabled(runing == false);
    if (runing) {
        ui->pushButton_xdb_start->setText("正在转换...");
    } else {
        ui->pushButton_xdb_start->setText("开始转换");
    }
}
