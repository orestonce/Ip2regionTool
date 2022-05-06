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
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_pushButton_input_db_clicked()
{
    QString input = QFileDialog::getOpenFileName(this, "",  "", "*.db");
    if (!input.isEmpty()) {
        ui->lineEdit_input_db->setText(input);
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
    req.DbFileName = db.toStdString();
    req.TxtFileName = txt.toStdString();
    req.Merge = ui->checkBox_DbToTxt_merge->isChecked();
    std::string errMsg = ConvertDbToTxt(req);
    if (errMsg.empty()) {
        QMessageBox::about(this, "成功", "转换成功!");
        return;
    }
    QMessageBox::warning(this, "错误", QString::fromStdString(errMsg));
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
        QMessageBox::about(this, "成功", "转换成功!");
        return;
    }
    QMessageBox::warning(this, "错误", QString::fromStdString(errMsg));
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

